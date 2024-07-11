package queue

import (
	"context"
	"encoding/json"
	"time"
)

var notificationEventID uint8 = 1
var notificationEventType = "sent"
var sendindInterval = 10 * time.Hour 

type notification struct {
	Rate   float64  `json:"rate"`
	Emails []string `json:"emails"`
	CreatedAt time.Time `json:"createdat"` 
	EventID uint8 `json:"eventid"` 
	EventType string `json:"eventtype"` 
	SendingTime time.Time `json:"sendingtime"` 
}

func (queue *Queue) Publish(ctx context.Context, rate float64, emails []string) error {
	n := notification{
		Rate:   rate,
		Emails: emails,
		CreatedAt: time.Unix(time.Now().Unix(), 0),
		EventID: notificationEventID,
		EventType: notificationEventType,
		SendingTime: time.Unix(time.Now().Add(sendindInterval).Unix(), 0),
	}

	payload, err := json.Marshal(n)

	if err != nil {
		return err
	}

	return queue.QueueConn.Publish(queue.subscriptionTopicName, payload)
}