package services

import (
	"context"
	"time"

	"github.com/RinaDish/currency-rates/tools"
)

const apiCallTimeout = 500 * time.Millisecond;

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
		ctx, cancel := context.WithTimeout(ctx, apiCallTimeout)
		defer cancel()

	return r.rateClient.GetDollarRate(ctx)
}