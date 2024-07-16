package queue

import (
	"github.com/RinaDish/currency-rates/tools"
)

type Broker interface {
	Publish(subj string, data []byte) error
	Drain() error 
}

type SubscriptionNotifierProducer struct {
	Broker Broker
	subscriptionTopicName string
	logger tools.Logger
}


func NewSubscriptionNotifierProducer(broker Broker, subscriptionTopicName string, logger tools.Logger) (*SubscriptionNotifierProducer) {
	return &SubscriptionNotifierProducer{
		Broker: broker,
		subscriptionTopicName: subscriptionTopicName,
		logger: logger.With("service", "queue"),
	}
}
