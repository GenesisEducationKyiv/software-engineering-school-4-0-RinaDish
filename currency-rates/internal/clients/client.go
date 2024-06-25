package clients

import (
	"context"
)

type RateClient interface {
	GetDollarRate(ctx context.Context) (float64, error) 
}

type Client struct {
	name string
	rateClient RateClient
}

func NewClient(name string, rateClient RateClient) *Client {
	return &Client{
		name: name,
		rateClient: rateClient,
	}
}

func (client *Client) GetDollarRate(ctx context.Context) (float64, error) {
	rate, err := client.rateClient.GetDollarRate(ctx)
	if err != nil {
		return 0.0, err
	}

	return rate, nil
}
