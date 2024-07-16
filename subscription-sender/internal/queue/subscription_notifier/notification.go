package subscriptionnotifier

import (
	"context"
	"encoding/json"

	"github.com/RinaDish/subscription-sender/internal/services"
	"github.com/nats-io/nats.go"
)

func (queue *SubscriptionNotifierConsumer) ConsumeSubscriptionEvent(ctx context.Context) error {
	_, err := queue.Broker.Subscribe(queue.subscriptionTopicName, func(msg *nats.Msg) {
        if err := queue.handleMessage(msg, ctx); err != nil {
			queue.logger.Error("Error handling queue message")
		}
    })

	if err != nil {
        return err
    }

	queue.logger.Infof("Listening for messages on %s topic...", queue.subscriptionTopicName)

	return nil
}

func  (queue *SubscriptionNotifierConsumer) handleMessage(msg *nats.Msg, ctx context.Context) error {
	var unmrshBody services.Notification
	err := json.Unmarshal(msg.Data, &unmrshBody)

	if err != nil {
		queue.logger.Error("Failed to unmarshal message")
		return err
	}

	queue.subscriptionService.NotifySubscribers(ctx, unmrshBody)
	queue.logger.Infof("Received a message: %s\n", string(msg.Data))

	return nil
}