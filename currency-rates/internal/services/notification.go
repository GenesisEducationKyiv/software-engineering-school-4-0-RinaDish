package services

import (
	"context"

	"github.com/RinaDish/currency-rates/tools"
)

type Email struct {
	ID       int    `json:"id" gorm:"id"`
	Email    string `json:"email" gorm:"email"`
	IsActive bool   `json:"is_active" gorm:"is_active"`
}

type SubscriptionDb interface {
	GetEmails(ctx context.Context) ([]Email, error)
}

type SubscriptionPublisher interface {
	Publish(ctx context.Context, rate float64, emails []string) error
}

type NotificationService struct {
	db           SubscriptionDb
	notification SubscriptionPublisher
	logger       tools.Logger
	rateClient   RateClient
}

func NewNotificationService(logger tools.Logger, d SubscriptionDb, s SubscriptionPublisher, r RateClient) NotificationService {
	return NotificationService{
		db:           d,
		notification: s,
		logger:       logger,
		rateClient:   r,
	}
}

func (service NotificationService) NotifySubscribers(ctx context.Context) error {
	rate, err := service.rateClient.GetDollarRate(ctx)

	if err != nil {
		service.logger.Error(err)
		return err
	}

	emails, err := service.db.GetEmails(ctx)

	if err != nil {
		service.logger.Error(err)
		return err
	}

	actualEmails := make([]string, 0, len(emails))
	for _, email := range emails {
		actualEmails = append(actualEmails, email.Email)
	}

	err = service.notification.Publish(ctx, rate, actualEmails)
	if err != nil {
		service.logger.Error(err)
		return err
	}

	service.logger.Infof("Message sent to subscription service...")

	return nil
}
