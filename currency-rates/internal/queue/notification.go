package queue

import (
	"context"
	"encoding/json"
	"time"
)

var notificationEventID uint8 = 1;
var notificationEventType = "sended";

type notification struct {
	Rate   float64  `json:"rate"`
	Emails []string `json:"emails"`
	Timestamp time.Time `json:"timestamp"` 
	EventID uint8 `json:"eventid"` 
	EventType string `json:"eventtype"` 
}

func (queue *Queue) Publish(ctx context.Context, rate float64, emails []string) error {
	n := notification{
		Rate:   rate,
		Emails: emails,
		Timestamp: time.Unix(time.Now().Unix(), 0),
		EventID: notificationEventID,
		EventType: notificationEventType,
	}

	payload, err := json.Marshal(n)

	if err != nil {
		return err
	}

	return queue.QueueConn.Publish(queue.subscriptionTopicName, payload)
}