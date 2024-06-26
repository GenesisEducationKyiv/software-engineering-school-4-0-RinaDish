package services

import (
	"context"

	"github.com/RinaDish/currency-rates/tools"
)

type RateClient interface {
	GetDollarRate(ctx context.Context) (float64, error) 
}

type Rate struct {
	logger tools.Logger
	rateClient RateClient
}

func NewRate(logger tools.Logger, rateClient RateClient) *Rate {
	return &Rate{
		logger: logger,
		rateClient: rateClient,
	}
}

func (r *Rate) GetDollarRate(ctx context.Context) (float64, error) {
	return r.rateClient.GetDollarRate(ctx)
}