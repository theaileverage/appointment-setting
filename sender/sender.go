package sender

import (
	"context"
	"fmt"

	"notification"

	"encore.dev/pubsub"
)

var notificationSub = pubsub.NewSubscription(notification.notificationTopic, "send-notification", pubsub.SubscriptionConfig{
    Handler: SendNotification,
})

//encore:api private
func SendNotification(ctx context.Context, n *notification.Notification) error {
    // In a real application, you would implement the actual sending logic here
    // For this example, we'll just print the notification details
    fmt.Printf("Sending notification: %+v\n", n)
    return nil
}