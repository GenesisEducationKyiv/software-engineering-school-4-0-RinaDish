package subscriptionnotifier

import (
	"context"
	"encoding/json"
)

func (queue *SubscriptionNotifierProducer) Publish(ctx context.Context, message interface{}) error {
	payload, err := json.Marshal(message)

	if err != nil {
		return err
	}

	return queue.Broker.Publish(queue.subscriptionTopicName, payload)
}