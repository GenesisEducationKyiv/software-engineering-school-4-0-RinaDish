package services

import (
	"context"

	"github.com/RinaDish/currency-rates/tools"
)

type CustomerService struct {
	db DB
	logger tools.Logger
}

func NewCustomerService(logger tools.Logger, db DB) CustomerService {
	return CustomerService{
		db: db,
		logger: logger,
	}
}

func (service CustomerService) CreateUser(ctx context.Context, email string) error {
    return service.db.SetUser(ctx, email)
}

func (service CustomerService) DeleteUser(ctx context.Context, email string) error {
    return service.db.DeleteUser(ctx, email) 
}
