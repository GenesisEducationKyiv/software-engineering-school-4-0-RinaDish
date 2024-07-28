package routers

import (
	"net/http"

	"github.com/RinaDish/currency-rates/tools"
	"github.com/go-chi/chi/v5"

	"github.com/RinaDish/currency-rates/internal/handlers"
)

type Metrics interface {
	MetricsHandler() http.Handler
    Middleware(next http.Handler) http.Handler
}

type Router struct {
	logger tools.Logger
	ratesHandler handlers.RateHandler
	subscriptionHandler handlers.SubscribeHandler
	metrics Metrics
}

func NewRouter(logger tools.Logger, ratesHandler handlers.RateHandler, subscriptionHandler handlers.SubscribeHandler, metrics Metrics) Router {
	return Router{
		logger: logger,
		ratesHandler: ratesHandler,
		subscriptionHandler: subscriptionHandler,
		metrics: metrics,
	}
}

func (router Router) GetRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(router.metrics.Middleware)

	r.Get("/rate", router.ratesHandler.GetCurrentRate)
	r.Post("/subscribe", router.subscriptionHandler.CreateSubscription)
	r.Post("/unsubscribe", router.subscriptionHandler.DeactivateSubscription)
	r.Get("/metrics", router.metrics.MetricsHandler().ServeHTTP)

	return r
}
