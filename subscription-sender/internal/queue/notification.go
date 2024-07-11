package queue

import (
	"context"

	"github.com/nats-io/nats.go"
)

func (queue *Queue) Subscribe(ctx context.Context) error {
	_, err := queue.QueueConn.Subscribe(queue.subscriptionTopicName, func(msg *nats.Msg) {
		queue.logger.Infof("Received a message: %s\n", string(msg.Data))
        queue.messagesService.GetMessage(msg, ctx)
    })

	if err != nil {
        return err
    }

	queue.logger.Infof("Listening for messages on 'subscription' topic...")

	return nil
}
