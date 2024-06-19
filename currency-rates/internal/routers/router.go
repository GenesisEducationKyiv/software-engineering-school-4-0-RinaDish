package routers

import (
	"go.uber.org/zap"
	"github.com/go-chi/chi/v5"

	"github.com/RinaDish/currency-rates/internal/handlers"
)

type Router struct {
	logger *zap.SugaredLogger
	ratesHandler handlers.RateHandler
	subscriptionHandler handlers.SubscribeHandler
}

func NewRouter(logger *zap.SugaredLogger, ratesHandler handlers.RateHandler, subscriptionHandler handlers.SubscribeHandler) Router {
	return Router{
		logger: logger,
		ratesHandler: ratesHandler,
		subscriptionHandler: subscriptionHandler,
	}
}

func (router Router) GetRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/rate", router.ratesHandler.GetCurrentRate)
	r.Post("/subscribe", router.subscriptionHandler.CreateSubscription)

	return r
}