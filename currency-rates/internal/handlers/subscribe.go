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
}

type SubscribeHandler struct {
	logger logger.Logger
	repo Db
}

func NewSubscribeHandler(logger logger.Logger, repo Db) SubscribeHandler {
	return SubscribeHandler{
		logger: logger,
		repo: repo,
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
	
	err = handler.repo.SetEmail(r.Context(), email)
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
  