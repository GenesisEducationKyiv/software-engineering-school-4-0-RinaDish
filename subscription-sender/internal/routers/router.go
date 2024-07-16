package routers

import (
	"github.com/RinaDish/subscription-sender/tools"
	"github.com/go-chi/chi/v5"

	"github.com/RinaDish/subscription-sender/internal/handlers"
)

type Router struct {
	logger tools.Logger
	healthCheckHandler handlers.HealthCheckHandler
}

func NewRouter(logger tools.Logger, healthCheckHandler handlers.HealthCheckHandler) Router {
	return Router{
		logger: logger,
		healthCheckHandler: healthCheckHandler,
	}
}

func (router Router) GetRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", router.healthCheckHandler.HealthCheck)

	return r
}