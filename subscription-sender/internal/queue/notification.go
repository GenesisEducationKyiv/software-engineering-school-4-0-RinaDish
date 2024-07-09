package queue

import (
	"context"
	"encoding/json"

	"github.com/RinaDish/subscription-sender/internal/services"
	"github.com/nats-io/nats.go"
)

func (queue *Queue) Subscribe(ctx context.Context, subscriptionService services.SubscriptionService) error {
	_, err := queue.QueueConn.Subscribe("subscription", func(msg *nats.Msg) {
        queue.getMessage(msg, ctx, subscriptionService)
    })

	if err != nil {
        return err
    }

	queue.logger.Infof("Listening for messages on 'subscription' topic...")

	return nil
}

func  (queue *Queue) getMessage(msg *nats.Msg, ctx context.Context, subscriptionService services.SubscriptionService) {
	var unmrshBody services.Notification
	_ = json.Unmarshal(msg.Data, &unmrshBody)

	subscriptionService.NotifySubscribers(ctx, unmrshBody)
	queue.logger.Infof("Received a message: %s\n", string(msg.Data))
}