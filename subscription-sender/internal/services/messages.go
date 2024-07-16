package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/RinaDish/subscription-sender/tools"
	"github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

type Message struct {
	ID    int    `json:"id" gorm:"id"`
	Rate   float64  `json:"rate" gorm:"rate"`
	Emails pq.StringArray  `json:"emails" gorm:"type:text[]"`
	CreatedAt time.Time `json:"createdat" gorm:"created_at"` 
	EventID uint8 `json:"eventid" gorm:"event_id"` 
	EventType string `json:"eventtype" gorm:"event_type"` 
	SendingTime time.Time `json:"sendingtime" gorm:"sending_time"`
	Sent bool `json:"sent" gorm:"sent"`
}

type MessagesDB interface {
	SetMessages(ctx context.Context, message Message) error
}

type MessagesService struct {
	db MessagesDB
	logger tools.Logger
}

func NewMessagesService(db MessagesDB, logger tools.Logger) (*MessagesService) {
	return &MessagesService{
		db:     db,
		logger: logger.With("service", "repository"),
	}
}

func  (messagesService *MessagesService) HandleMessage(msg *nats.Msg, ctx context.Context) {
	var unmrshMsg Message
	_ = json.Unmarshal(msg.Data, &unmrshMsg)

	unmrshMsg.Sent = false;

	err := messagesService.db.SetMessages(ctx, unmrshMsg)

	if err != nil {
		messagesService.logger.Error(err)
	}
}