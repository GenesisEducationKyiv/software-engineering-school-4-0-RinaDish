package services

import (
	"context"
	"fmt"

	"github.com/RinaDish/currency-rates/internal/repo"
	"go.uber.org/zap"
)

type SubscriptionDb interface {
	GetEmails(ctx context.Context) ([]repo.Email, error)
}

type SubscriptionSender interface {
	Send(to, body string)
}

type SubscriptionService struct {
	db SubscriptionDb
	sender SubscriptionSender
	l *zap.SugaredLogger
	r RateClient
}

func NewSubscriptionService(l *zap.SugaredLogger, d SubscriptionDb, s SubscriptionSender, r RateClient) SubscriptionService{

	return SubscriptionService{
		db: d,
		sender: s,
		l: l,
		r: r,
	}
}

func (s SubscriptionService) NotifySubscribers(ctx context.Context){
	rate, err := s.r.GetDollarRate(ctx)

	if err != nil {
		s.l.Error(err)
		return
	}

	emails, err := s.db.GetEmails(ctx)

	if err != nil {
		s.l.Error(err)
		return
	}

	for _, email := range emails {
		s.sender.Send(email.Email, fmt.Sprintf("%f", rate))
	}
}