package services

import (
	"context"
	"strconv"

	"github.com/RinaDish/subscription-sender/tools"
)

type Notification struct {
	Rate float64    `json:"rate"`
	Emails []string `json:"emails"`
}

type SubscriptionSender interface {
	Send(to, body string)
}

type SubscriptionRepository interface {
	GetMessages(ctx context.Context) ([]Message, error)
	UpdateMessages(ctx context.Context, msgID int) error
}

type SubscriptionService struct {
	sender SubscriptionSender
	repo SubscriptionRepository
	logger tools.Logger
}

func NewSubscriptionService(logger tools.Logger, sender SubscriptionSender, repo SubscriptionRepository) SubscriptionService{
	return SubscriptionService{
		sender: sender,
		repo: repo,
		logger: logger,
	}
}

func (service SubscriptionService) NotifySubscribers(ctx context.Context) error {
	msgs, err := service.repo.GetMessages(ctx) 

	if err != nil {
		return err
	}

	for _, msg := range msgs {
		for _, email := range msg.Emails {
			service.sender.Send(email, strconv.FormatFloat(msg.Rate, 'g', -1, 64))
		}

		err =  service.repo.UpdateMessages(ctx, msg.ID)

		if err != nil {
			return err
		}
	}

	return nil
}