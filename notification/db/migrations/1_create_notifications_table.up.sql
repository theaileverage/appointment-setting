CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message TEXT NOT NULL,
    channel TEXT NOT NULL,
    recipient TEXT NOT NULL,
    send_at TIMESTAMP WITH TIME ZONE NOT NULL,
    sent BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_notifications_send_at ON notifications(send_at);
CREATE INDEX idx_notifications_sent ON notifications(sent);
