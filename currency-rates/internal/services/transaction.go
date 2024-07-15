package services

import (
	"context"

	"github.com/RinaDish/currency-rates/tools"
)

type DB interface {
	SetEmail(ctx context.Context, email string) error
	DeactivateEmail(ctx context.Context, email string) error
	SetUser(ctx context.Context, email string) error
	DeleteUser(ctx context.Context, email string) error
}

type Transaction struct {
	customerService     CustomerService
	subscriptionService SubscriptionService
	logger              tools.Logger
}

func NewTransactionService(logger tools.Logger, customerService CustomerService, subscriptionService SubscriptionService) Transaction {
	return Transaction{
		logger:              logger,
		subscriptionService: subscriptionService,
		customerService:     customerService,
	}
}

func (transaction Transaction) ExecuteSubscription(ctx context.Context, email string) error {
	transaction.logger.Info("Starting SAGA transaction")

	
	if err := transaction.createUser(ctx, email); err != nil {
		return err
	}

	
	if err := transaction.createSubscription(ctx, email); err != nil {
		err = transaction.compensateCreateUser(ctx, email)
		return err
	}

	transaction.logger.Info("SAGA transaction completed successfully")

	return nil
}

func (transaction Transaction) createUser(ctx context.Context, email string) error {
	transaction.logger.Info("Creating user")
	return transaction.customerService.CreateUser(ctx, email)
}

func (transaction Transaction) compensateCreateUser(ctx context.Context, email string) error {
	transaction.logger.Info("Compensating create user")
	return transaction.customerService.DeleteUser(ctx, email)
}

func (transaction Transaction) createSubscription(ctx context.Context, email string) error {
	transaction.logger.Info("Creating subscription")
	return transaction.subscriptionService.CreateSubscription(ctx, email)
}
