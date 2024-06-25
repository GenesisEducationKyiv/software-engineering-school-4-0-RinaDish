package clients

import (
	"context"
)

type Chain interface {
    RateClient
    SetNext(Chain)
}

type BaseChain struct {
    rateClient RateClient
    next Chain
}

func NewBaseChain(rateClient RateClient) *BaseChain {
    return &BaseChain{
        rateClient: rateClient,
    }
}

func (chain *BaseChain) SetNext(next Chain) {
    chain.next = next
}

func (chain *BaseChain) GetDollarRate(ctx context.Context) (float64, error) {
	rate, err := chain.rateClient.GetDollarRate(ctx)

	if err != nil {
		next := chain.next
		if next == nil {
			return 0.0, err
		}

		return next.GetDollarRate(ctx)
	}

	return rate, nil
}