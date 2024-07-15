package handlers

import (
	"context"
	"net/http"
	"regexp"

	"github.com/RinaDish/currency-rates/tools"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type Db interface {
	SetEmail(ctx context.Context, email string) error
	DeactivateEmail(ctx context.Context, email string) error
}

type Transaction interface {
	ExecuteSubscription(ctx context.Context, email string) error
}

type SubscriptionService interface {
	DeactivateSubscription(ctx context.Context, email string) error
}

type SubscribeHandler struct {
	logger              tools.Logger
	repo                Db
	transaction         Transaction
	subscriptionService SubscriptionService
}

func NewSubscribeHandler(logger tools.Logger, repo Db, transaction Transaction, subscriptionService SubscriptionService) SubscribeHandler {
	return SubscribeHandler{
		logger:              logger,
		repo:                repo,
		transaction:         transaction,
		subscriptionService: subscriptionService,
	}
}

func (handler SubscribeHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form parse error", http.StatusConflict)
		return
	}

	formData := r.Form

	email := formData.Get("email")

	if !isValidEmail(email) {
		http.Error(w, "Invalid email", http.StatusConflict)
		handler.logger.Info("Invalid email")

		return
	}

	err = handler.transaction.ExecuteSubscription(r.Context(), email)

	responseStatus := http.StatusOK
	if err != nil {
		responseStatus = http.StatusConflict
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
}

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func (handler SubscribeHandler) DeactivateSubscription(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form parse error", http.StatusBadRequest)
		return
	}

	formData := r.Form

	email := formData.Get("email")

	err = handler.subscriptionService.DeactivateSubscription(r.Context(), email)
	responseStatus := http.StatusOK
	if err != nil {
		responseStatus = http.StatusNotFound
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
}
