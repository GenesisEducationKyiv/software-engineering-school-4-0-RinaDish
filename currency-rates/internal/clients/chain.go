package clients

import (
	"context"
	"time"
)

const ApiCallTimeout = 500 * time.Millisecond;

type RateClient interface {
	GetDollarRate(ctx context.Context) (float64, error) 
}

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
		if chain.next == nil {
			return 0.0, err
		}

		return chain.next.GetDollarRate(ctx)
	}

	return rate, nil
}