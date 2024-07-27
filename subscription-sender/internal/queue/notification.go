package queue

import (
	"context"

	"github.com/nats-io/nats.go"
)

func (queue *SubscriptionNotifierConsumer) ConsumeSubscriptionEvent(ctx context.Context) error {
	_, err := queue.Broker.Subscribe(queue.subscriptionTopicName, func(msg *nats.Msg) {
        queue.messagesService.HandleMessage(msg.Data, ctx)
    })

	if err != nil {
        return err
    }

	queue.logger.Infof("Listening for messages on %s topic...", queue.subscriptionTopicName)

	return nil
}
