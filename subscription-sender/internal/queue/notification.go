package queue

import (
	"context"
	"encoding/json"

	"github.com/RinaDish/subscription-sender/internal/services"
	"github.com/nats-io/nats.go"
)

func (queue *Queue) Subscribe(ctx context.Context) error {
	_, err := queue.QueueConn.Subscribe(queue.subscriptionTopicName, func(msg *nats.Msg) {
        queue.getMessage(msg, ctx)
    })

	if err != nil {
        return err
    }

	queue.logger.Infof("Listening for messages on 'subscription' topic...")

	return nil
}

func  (queue *Queue) getMessage(msg *nats.Msg, ctx context.Context) {
	var unmrshBody services.Notification
	_ = json.Unmarshal(msg.Data, &unmrshBody)

	queue.subscriptionService.NotifySubscribers(ctx, unmrshBody)
	queue.logger.Infof("Received a message: %s\n", string(msg.Data))
}