package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/RinaDish/currency-rates/tools"
)

type notificationRequest struct {
	Rate   float64  `json:"rate"`
	Emails []string `json:"emails"`
}

type NotificationClient struct {
	logger tools.Logger
	client *http.Client
	url string
}

func NewNotificationClient(notificationServiceURL string, logger tools.Logger) NotificationClient {
	return NotificationClient{
		logger: logger.With("client", "Notification"),
		client: http.DefaultClient,
		url: notificationServiceURL,
	}
}

func (notificationClient NotificationClient) Send(ctx context.Context, rate float64, emails []string) error {
	n := notificationRequest{
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

	if res.StatusCode != http.StatusOK {
		return errors.New("notification failed")
	}

	return nil
}
