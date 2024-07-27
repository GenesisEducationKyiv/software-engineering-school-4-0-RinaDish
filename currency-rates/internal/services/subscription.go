package services

import (
	"context"
	"time"

	"github.com/RinaDish/currency-rates/tools"
)

var notificationEventID uint8 = 1;
var notificationEventType = "sent";
var sendindInterval = 10 * time.Hour 

type notification struct {
	Rate   float64  `json:"rate"`
	Emails []string `json:"emails"`
	Timestamp time.Time `json:"timestamp"` 
	EventID uint8 `json:"eventid"` 
	EventType string `json:"eventtype"` 
	SendingTime time.Time `json:"sendingtime"` 
}

type Email struct {
	ID    int    `json:"id" gorm:"id"`
	Email string `json:"email" gorm:"email"`
	IsActive bool `json:"is_active" gorm:"is_active"`
}

type SubscriptionDb interface {
	GetAllActiveEmails(ctx context.Context) ([]Email, error)
}

type SubscriptionPublisher interface {
	Publish(ctx context.Context, message interface{}) error
}

type SubscriptionService struct {
	db SubscriptionDb
	notification SubscriptionPublisher
	logger tools.Logger
	rateClient RateClient
}

func NewSubscriptionService(logger tools.Logger, d SubscriptionDb, s SubscriptionPublisher, r RateClient) SubscriptionService{
	return SubscriptionService{
		db: d,
		notification: s,
		logger: logger,
		rateClient: r,
	}
}

func (service SubscriptionService) NotifySubscribers(ctx context.Context) error {
	rate, err := service.rateClient.GetDollarRate(ctx)

	if err != nil {
		service.logger.Error(err)
		return err
	}

	emails, err := service.db.GetAllActiveEmails(ctx)

	if err != nil {
		service.logger.Error(err)
		return err
	}

	actualEmails := make([]string, 0, len(emails))
	for _, email := range emails {
		actualEmails = append(actualEmails, email.Email)
	}

	n := notification{
		Rate:   rate,
		Emails: actualEmails,
		Timestamp: time.Unix(time.Now().Unix(), 0),
		EventID: notificationEventID,
		EventType: notificationEventType,
		SendingTime: time.Unix(time.Now().Add(sendindInterval).Unix(), 0),
	}

	err = service.notification.Publish(ctx, n)
	if err != nil {
		service.logger.Error(err)
		return err
	}

	service.logger.Infof("Messege sent to subscription service...")

	return nil
}