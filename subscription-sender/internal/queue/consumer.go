package queue

import (
	"context"

	"github.com/nats-io/nats.go"

	"github.com/RinaDish/subscription-sender/tools"
)

type Broker interface {
	Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error)
	Drain() error 
}

type notificationService interface {
	HandleMessage(msg []byte, ctx context.Context)
}

type SubscriptionNotifierConsumer struct {
	Broker Broker
	messagesService notificationService
	subscriptionTopicName string
	logger tools.Logger
}


func NewSubscriptionNotifierConsumer(broker Broker, subscriptionTopicName string, messagesService notificationService, logger tools.Logger) (*SubscriptionNotifierConsumer) {
	return &SubscriptionNotifierConsumer{
		Broker: broker,
		messagesService: messagesService,
		subscriptionTopicName: subscriptionTopicName,
		logger: logger.With("service", "queue"),
	}
}