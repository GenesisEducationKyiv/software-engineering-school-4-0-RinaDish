package queue

import (
	"github.com/nats-io/nats.go"

	"github.com/RinaDish/subscription-sender/internal/services"
	"github.com/RinaDish/subscription-sender/tools"
)

type Queue struct {
	QueueConn *nats.Conn
	subscriptionService services.SubscriptionService
	subscriptionTopicName string
	logger tools.Logger
}


func NewQueue(nats *nats.Conn,subscriptionTopicName string, subscriptionService services.SubscriptionService, logger tools.Logger) (*Queue) {
	return &Queue{
		QueueConn:     nats,
		subscriptionService: subscriptionService,
		subscriptionTopicName: subscriptionTopicName,
		logger: logger.With("service", "queue"),
	}
}
