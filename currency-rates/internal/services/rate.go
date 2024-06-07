package services

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type RateClient interface {
	GetDollarRate(ctx context.Context) (float64, error) 
}

type Rate struct {
	l *zap.SugaredLogger
	rateClients []RateClient
}

func NewRate(l *zap.SugaredLogger, r []RateClient) Rate {
	return Rate{
		l: l,
		rateClients: r,
	}
}

func (r Rate) GetDollarRate(ctx context.Context) (float64, error) {
	var rate float64
	var err error
	for _, c := range r.rateClients {
		ctx, _ := context.WithTimeout(ctx, 500 * time.Millisecond)
		rate, err = c.GetDollarRate(ctx)
		if err == nil { break }
	}

	return rate, err
}