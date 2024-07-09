package queue

import (
	"github.com/nats-io/nats.go"

	"github.com/RinaDish/currency-rates/tools"
)

type Queue struct {
	QueueConn *nats.Conn
	subscriptionTopicName string
	logger tools.Logger
}


func NewQueue(nats *nats.Conn, subscriptionTopicName string, logger tools.Logger) (*Queue) {
	return &Queue{
		QueueConn: nats,
		subscriptionTopicName: subscriptionTopicName,
		logger: logger.With("service", "queue"),
	}
}
