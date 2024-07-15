package services

import (
	"context"
	"strconv"

	"github.com/RinaDish/subscription-sender/tools"
)

type Notification struct {
	Rate   float64  `json:"rate"`
	Emails []string `json:"emails"`
}

type SubscriptionSender interface {
	Send(to, body string)
}

type SubscriptionRepository interface {
	GetMessages(ctx context.Context) ([]Message, error)
	UpdateMessages(ctx context.Context, message Message) error
}

type NotificationService struct {
	sender SubscriptionSender
	repo   SubscriptionRepository
	logger tools.Logger
}

func NewNotificationService(logger tools.Logger, sender SubscriptionSender, repo SubscriptionRepository) NotificationService {
	return NotificationService{
		sender: sender,
		repo:   repo,
		logger: logger,
	}
}

func (service NotificationService) NotifySubscribers(ctx context.Context) error {
	msgs, err := service.repo.GetMessages(ctx)

	if err != nil {
		return err
	}

	for _, msg := range msgs {
		for _, email := range msg.Emails {
			service.sender.Send(email, strconv.FormatFloat(msg.Rate, 'g', -1, 64))
		}

		err = service.repo.UpdateMessages(ctx, msg)

		if err != nil {
			return err
		}
	}

	return nil
}
