package services

import (
	"context"
	"time"

	"go.uber.org/zap"
)

const apiCallTimeout = 500 * time.Millisecond;

type RateClient interface {
	GetDollarRate(ctx context.Context) (float64, error) 
}

type Rate struct {
	logger *zap.SugaredLogger
	rateClients []RateClient
}

func NewRate(logger *zap.SugaredLogger, rateClients []RateClient) Rate {
	return Rate{
		logger: logger,
		rateClients: rateClients,
	}
}

func (r Rate) GetDollarRate(ctx context.Context) (float64, error) {
	var rate float64
	var err error

	for _, c := range r.rateClients {
		ctx, cancel := context.WithTimeout(ctx, apiCallTimeout)
		defer cancel()

		rate, err = c.GetDollarRate(ctx)
		if err == nil { break }
	}

	return rate, err
}