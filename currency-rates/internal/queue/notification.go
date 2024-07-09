package queue

import (
	"context"
	"encoding/json"
)

type Notification struct {
	Rate   float64  `json:"rate"`
	Emails []string `json:"emails"`
}

func (queue *Queue) Publish(ctx context.Context, rate float64, emails []string) error {
	n := Notification{
		Rate:   rate,
		Emails: emails,
	}

	payload, err := json.Marshal(n)

	if err != nil {
		return err
	}

	return queue.QueueConn.Publish("subscription", payload)
}