package services

import (
	"context"

	"github.com/RinaDish/currency-rates/tools"
)

type SubscriptionService struct {
	db DB
	logger tools.Logger
}

func NewSubscriptionService(logger tools.Logger, db DB) SubscriptionService {
	return SubscriptionService{
		db: db,
		logger: logger,
	}
}

func (service SubscriptionService) CreateSubscription(ctx context.Context, email string) error {
    return service.db.SetEmail(ctx, email) 
}

func (service SubscriptionService) DeactivateSubscription(ctx context.Context, email string) error {
    return service.db.DeactivateEmail(ctx, email) 
}
