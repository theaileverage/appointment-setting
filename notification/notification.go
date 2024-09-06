package notification

import (
	"context"
	"errors"
	"time"

	"encore.dev/cron"
	"encore.dev/pubsub"
	"encore.dev/storage/sqldb"
)

// Define the database
var db = sqldb.NewDatabase("notifications", sqldb.DatabaseConfig{
	Migrations: "./db/migrations",
})

type Channel string

const (
	ChannelWhatsApp Channel = "whatsapp"
	ChannelEmail    Channel = "email"
	ChannelTelegram Channel = "telegram"
	ChannelSlack    Channel = "slack"
	ChannelDiscord  Channel = "discord"
)

type Notification struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Channel   Channel   `json:"channel"`
	Recipient string    `json:"recipient"`
	SendAt    time.Time `json:"send_at"`
	Sent      bool      `json:"sent"`
}

type CreateParams struct {
	Message   string    `json:"message"`
	Channel   Channel   `json:"channel"`
	Recipient string    `json:"recipient"`
	SendAt    time.Time `json:"send_at"`
}

type ListResponse struct {
	Notifications []*Notification `json:"notifications"`
}

type UpdateParams struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Channel   Channel   `json:"channel"`
	Recipient string    `json:"recipient"`
	SendAt    time.Time `json:"send_at"`
}

//encore:api public method=POST path=/notifications
func Create(ctx context.Context, params *CreateParams) (*Notification, error) {
	notification := &Notification{
		Message:   params.Message,
		Channel:   params.Channel,
		Recipient: params.Recipient,
		SendAt:    params.SendAt,
		Sent:      false,
	}

	err := db.QueryRow(ctx, `
		INSERT INTO notifications (message, channel, recipient, send_at, sent)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, notification.Message, notification.Channel, notification.Recipient, notification.SendAt, notification.Sent).Scan(&notification.ID)

	if err != nil {
		return nil, err
	}

	return notification, nil
}

//encore:api public method=GET path=/notifications
func List(ctx context.Context) (*ListResponse, error) {
	rows, err := db.Query(ctx, `
		SELECT id, message, channel, recipient, send_at, sent
		FROM notifications
		WHERE sent = false
		ORDER BY send_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*Notification
	for rows.Next() {
		n := &Notification{}
		if err := rows.Scan(&n.ID, &n.Message, &n.Channel, &n.Recipient, &n.SendAt, &n.Sent); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &ListResponse{Notifications: notifications}, nil
}

//encore:api public method=PUT path=/notifications/:id
func Update(ctx context.Context, id string, params *UpdateParams) (*Notification, error) {
	result, err := db.Exec(ctx, `
		UPDATE notifications
		SET message = $1, channel = $2, recipient = $3, send_at = $4
		WHERE id = $5 AND sent = false
	`, params.Message, params.Channel, params.Recipient, params.SendAt, id)

	if err != nil {
		return nil, err
	}

	affected := result.RowsAffected()
	if affected == 0 {
		return nil, errors.New("notification not found or already sent")
	}

	// Fetch the updated notification
	var notification Notification
	err = db.QueryRow(ctx, `
		SELECT id, message, channel, recipient, send_at, sent
		FROM notifications
		WHERE id = $1
	`, id).Scan(&notification.ID, &notification.Message, &notification.Channel, &notification.Recipient, &notification.SendAt, &notification.Sent)

	if err != nil {
		return nil, err
	}

	return &notification, nil
}

//encore:api public method=DELETE path=/notifications/:id
func Delete(ctx context.Context, id string) error {
	result, err := db.Exec(ctx, "DELETE FROM notifications WHERE id = $1", id)
	if err != nil {
		return err
	}
	affected := result.RowsAffected()
	if affected == 0 {
		return errors.New("notification not found")
	}
	return nil
}

var notificationTopic = pubsub.NewTopic[*Notification]("notifications", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})

//encore:api public method=POST path=/notifications/:id/send
func Send(ctx context.Context, id string) error {
	result, err := db.Exec(ctx, `
		UPDATE notifications
		SET sent = true
		WHERE id = $1 AND sent = false
	`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("notification not found or already sent")
	}

	// Fetch the notification
	var n Notification
	err = db.QueryRow(ctx, `
		SELECT id, message, channel, recipient, send_at, sent
		FROM notifications
		WHERE id = $1
	`, id).Scan(&n.ID, &n.Message, &n.Channel, &n.Recipient, &n.SendAt, &n.Sent)
	if err != nil {
		return err
	}

	// Publish the notification to be sent
	_, err = notificationTopic.Publish(ctx, &n)
	return err
}

var _ = cron.NewJob("send-notifications", cron.JobConfig{
	Title:    "Send scheduled notifications",
	Endpoint: SendScheduledNotifications,
	Every:    1 * cron.Minute,
})

//encore:api private
func SendScheduledNotifications(ctx context.Context) error {
	rows, err := db.Query(ctx, `
		SELECT id
		FROM notifications
		WHERE sent = false AND send_at <= NOW()
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return err
		}
		if err := Send(ctx, id); err != nil {
			return err
		}
	}
	return rows.Err()
}
