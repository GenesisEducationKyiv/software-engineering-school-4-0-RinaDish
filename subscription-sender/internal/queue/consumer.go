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
	HandleMessage(msg *nats.Msg, ctx context.Context)
}

type Queue struct {
	Broker Broker
	messagesService notificationService
	subscriptionTopicName string
	logger tools.Logger
}


func NewQueue(broker Broker, subscriptionTopicName string, messagesService notificationService, logger tools.Logger) (*Queue) {
	return &Queue{
		Broker: broker,
		messagesService: messagesService,
		subscriptionTopicName: subscriptionTopicName,
		logger: logger.With("service", "queue"),
	}
}