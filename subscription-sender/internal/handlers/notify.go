package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/RinaDish/subscription-sender/tools"
	"github.com/RinaDish/subscription-sender/internal/services"
)

type Notification struct {
	Rate float64    `json:"rate"`
	Emails []string `json:"emails"`
}

type NotifyService interface {
	NotifySubscribers(ctx context.Context, notification services.Notification)
}

type NotifyHandler struct {
	logger tools.Logger
	subscriptionService NotifyService
}

func NewNotifyHandler(logger tools.Logger, subscriptionService NotifyService) NotifyHandler {
	return NotifyHandler{
		logger: logger,
		subscriptionService: subscriptionService,
	}
}

func (handler NotifyHandler) NotifySubscribers(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	var unmrshBody Notification
	err = json.Unmarshal(body, &unmrshBody)
	
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	handler.logger.Info(unmrshBody)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	handler.subscriptionService.NotifySubscribers(context.Background(), services.Notification{Rate: unmrshBody.Rate, Emails: unmrshBody.Emails})
}
  