package queue

import (
	"context"

	"github.com/nats-io/nats.go"

	"github.com/RinaDish/subscription-sender/internal/services"
	"github.com/RinaDish/subscription-sender/tools"
)

type Broker interface {
	Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error)
	Drain() error 
}

type NotificationService interface {
	NotifySubscribers(ctx context.Context, notification services.Notification )
}

type Queue struct {
	Broker Broker
	subscriptionService NotificationService
	subscriptionTopicName string
	logger tools.Logger
}


func NewQueue(broker Broker, subscriptionTopicName string, subscriptionService NotificationService, logger tools.Logger) (*Queue) {
	return &Queue{
		Broker: broker,
		subscriptionService: subscriptionService,
		subscriptionTopicName: subscriptionTopicName,
		logger: logger.With("service", "queue"),
	}
}
