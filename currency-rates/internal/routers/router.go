package routers

import (
	"github.com/RinaDish/currency-rates/tools"
	"github.com/go-chi/chi/v5"

	"github.com/RinaDish/currency-rates/internal/handlers"
)

type Router struct {
	logger tools.Logger
	ratesHandler handlers.RateHandler
	subscriptionHandler handlers.SubscribeHandler
}

func NewRouter(logger tools.Logger, ratesHandler handlers.RateHandler, subscriptionHandler handlers.SubscribeHandler) Router {
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