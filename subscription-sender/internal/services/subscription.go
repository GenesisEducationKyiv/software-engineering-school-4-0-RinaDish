package services

import (
	"context"
	"fmt"

	"github.com/RinaDish/subscription-sender/tools"
)

type Notification struct {
	Rate float64    `json:"rate"`
	Emails []string `json:"emails"`
}

type SubscriptionSender interface {
	Send(to, body string)
}

type SubscriptionService struct {
	sender SubscriptionSender
	logger tools.Logger
}

func NewSubscriptionService(logger tools.Logger, sender SubscriptionSender) SubscriptionService{
	return SubscriptionService{
		sender: sender,
		logger: logger,
	}
}

func (service SubscriptionService) NotifySubscribers(ctx context.Context, notification Notification ){
	for _, email := range notification.Emails {
		service.sender.Send(email, fmt.Sprintf("%f", notification.Rate))
	}
}