package services

import (
	"context"
	"fmt"

	"github.com/RinaDish/currency-rates/internal/repo"
	"github.com/RinaDish/currency-rates/tools"
)

type SubscriptionDb interface {
	GetEmails(ctx context.Context) ([]repo.Email, error)
}

type SubscriptionSender interface {
	Send(to, body string)
}

type SubscriptionService struct {
	db SubscriptionDb
	sender SubscriptionSender
	logger logger.Logger
	rateClient RateClient
}

func NewSubscriptionService(logger logger.Logger, d SubscriptionDb, s SubscriptionSender, r RateClient) SubscriptionService{
	return SubscriptionService{
		db: d,
		sender: s,
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

	for _, email := range emails {
		service.sender.Send(email.Email, fmt.Sprintf("%f", rate))
	}
}