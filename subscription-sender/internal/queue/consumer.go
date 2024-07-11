package queue

import (
	"context"

	"github.com/nats-io/nats.go"

	"github.com/RinaDish/subscription-sender/tools"
)

type notificationService interface {
	GetMessage(msg *nats.Msg, ctx context.Context)
}

type Queue struct {
	QueueConn *nats.Conn
	messagesService notificationService
	subscriptionTopicName string
	logger tools.Logger
}


func NewQueue(nats *nats.Conn,subscriptionTopicName string, messagesService notificationService, logger tools.Logger) (*Queue) {
	return &Queue{
		QueueConn:     nats,
		messagesService: messagesService,
		subscriptionTopicName: subscriptionTopicName,
		logger: logger.With("service", "queue"),
	}
}
