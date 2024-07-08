package handlers

import (
	"net/http"

	"github.com/RinaDish/subscription-sender/tools"
)

type HealthCheckHandler struct {
	logger tools.Logger
}

func NewHealthCheckHandler(logger tools.Logger) HealthCheckHandler {
	return HealthCheckHandler{
		logger: logger,
	}
}

func (handler HealthCheckHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
  