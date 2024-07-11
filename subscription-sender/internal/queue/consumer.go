package queue

import (
	"context"

	"github.com/nats-io/nats.go"

	"github.com/RinaDish/subscription-sender/internal/services"
	"github.com/RinaDish/subscription-sender/tools"
)

type NotificationService interface {
	NotifySubscribers(ctx context.Context, notification services.Notification )
}

type Queue struct {
	QueueConn *nats.Conn
	subscriptionService NotificationService
	subscriptionTopicName string
	logger tools.Logger
}


func NewQueue(nats *nats.Conn, subscriptionTopicName string, subscriptionService NotificationService, logger tools.Logger) (*Queue) {
	return &Queue{
		QueueConn:     nats,
		subscriptionService: subscriptionService,
		subscriptionTopicName: subscriptionTopicName,
		logger: logger.With("service", "queue"),
	}
}
