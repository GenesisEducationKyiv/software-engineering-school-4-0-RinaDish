package queue

import (
	"github.com/RinaDish/currency-rates/tools"
)

type Broker interface {
	Publish(subj string, data []byte) error
	Drain() error 
}

type Queue struct {
	Broker Broker
	subscriptionTopicName string
	logger tools.Logger
}


func NewQueue(broker Broker, subscriptionTopicName string, logger tools.Logger) (*Queue) {
	return &Queue{
		Broker: broker,
		subscriptionTopicName: subscriptionTopicName,
		logger: logger.With("service", "queue"),
	}
}
