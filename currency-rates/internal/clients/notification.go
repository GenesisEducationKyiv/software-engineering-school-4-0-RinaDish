package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/RinaDish/currency-rates/tools"
)

type Notification struct {
	Rate   float64  `json:"rate"`
	Emails []string `json:"emails"`
}

type NotificationClient struct {
	logger tools.Logger
	client *http.Client
	url string
}

func NewNotificationClient(notificationServiceUrl string, logger tools.Logger) NotificationClient {
	client := http.DefaultClient

	return NotificationClient{
		logger: logger.With("client", "Notification"),
		client: client,
		url: notificationServiceUrl,
	}
}

func (notificationClient NotificationClient) Send(ctx context.Context, rate float64, emails []string) error {
	n := Notification{
		Rate:   rate,
		Emails: emails,
	}

	payload, err := json.Marshal(n)

	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, notificationClient.url, bytes.NewBuffer(payload))

	if err != nil {
		return err
	}

	res, err := notificationClient.client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}
