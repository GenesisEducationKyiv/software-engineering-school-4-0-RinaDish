package queue

import (
	"github.com/nats-io/nats.go"

	"github.com/RinaDish/subscription-sender/internal/services"
	"github.com/RinaDish/subscription-sender/tools"
)

type Queue struct {
	QueueConn *nats.Conn
	messagesService *services.MessagesService
	subscriptionTopicName string
	logger tools.Logger
}


func NewQueue(nats *nats.Conn,subscriptionTopicName string, messagesService *services.MessagesService, logger tools.Logger) (*Queue) {
	return &Queue{
		QueueConn:     nats,
		messagesService: messagesService,
		subscriptionTopicName: subscriptionTopicName,
		logger: logger.With("service", "queue"),
	}
}
