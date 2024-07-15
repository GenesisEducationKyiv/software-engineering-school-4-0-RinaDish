package queue

import (
	"context"
	"encoding/json"
)

func (queue *Queue) Publish(ctx context.Context, message interface{}) error {
	payload, err := json.Marshal(message)

	if err != nil {
		return err
	}

	return queue.QueueConn.Publish(queue.subscriptionTopicName, payload)
}