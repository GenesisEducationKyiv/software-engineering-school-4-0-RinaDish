package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/RinaDish/subscription-sender/tools"
)

type Message struct {
	ID int `json:"id"`
	Rate float64 `json:"rate"`
	Emails []string `json:"emails"`
	CreatedAt time.Time `json:"createdat"`
	EventID uint8 `json:"eventid"`
	EventType string `json:"eventtype"`
	SendingTime time.Time `json:"sendingtime"`
	Sent bool `json:"sent"`
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

func  (messagesService *MessagesService) HandleMessage(msg []byte, ctx context.Context) {
	var unmrshMsg Message
	_ = json.Unmarshal(msg, &unmrshMsg)

	unmrshMsg.Sent = false
	messagesService.logger.Info("Received message: ", unmrshMsg)
	err := messagesService.db.SetMessages(ctx, unmrshMsg)

	if err != nil {
		messagesService.logger.Error(err)
	}
}
