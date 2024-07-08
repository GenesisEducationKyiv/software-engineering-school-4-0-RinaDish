package services

import (
	"context"
	"fmt"

	"github.com/RinaDish/currency-rates/tools"
)

type Email struct {
	ID    int    `json:"id" gorm:"id"`
	Email string `json:"email" gorm:"email"`
}

type SubscriptionDb interface {
	GetEmails(ctx context.Context) ([]Email, error)
}

type SubscriptionPublisher interface {
	Publish(ctx context.Context, rate float64, emails []string) error
}

type SubscriptionService struct {
	db SubscriptionDb
	notification SubscriptionPublisher
	logger tools.Logger
	rateClient RateClient
}

func NewSubscriptionService(logger tools.Logger, d SubscriptionDb, s SubscriptionPublisher, r RateClient) SubscriptionService{
	return SubscriptionService{
		db: d,
		notification: s,
		logger: logger,
		rateClient: r,
	}
}

func (service SubscriptionService) NotifySubscribers(ctx context.Context){
	rate, err := service.rateClient.GetDollarRate(ctx)

	if err != nil {
		service.logger.Error(err)
		return
	}

	emails, err := service.db.GetEmails(ctx)

	if err != nil {
		service.logger.Error(err)
		return
	}

	actualEmails := make([]string, 0, len(emails))
	for _, email := range emails {
		actualEmails = append(actualEmails, email.Email)
	}

	fmt.Println(actualEmails)

	err = service.notification.Publish(ctx, rate, actualEmails)
	if err != nil {
		service.logger.Error(err)
		return
	}
}