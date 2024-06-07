package handlers

import (
	"context"
	"net/http"
	"regexp"

	"go.uber.org/zap"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type Db interface {
	SetEmail(ctx context.Context, email string) error
}

type SubscribeHandler struct {
	l *zap.SugaredLogger
	r Db
}

func NewSubscribeHandler(l *zap.SugaredLogger, r Db) SubscribeHandler {

	return SubscribeHandler{
		l: l,
		r: r,
	}
}

func (h SubscribeHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form parse error", http.StatusConflict)
		return
	}

	formData := r.Form

	email := formData.Get("email")

	if !isValidEmail(email) {
		http.Error(w, "Invalid email", http.StatusConflict)
		return
	}
	
	err = h.r.SetEmail(r.Context(), email)
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
  