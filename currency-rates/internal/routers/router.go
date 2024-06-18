package routers

import (
	"go.uber.org/zap"
	"github.com/go-chi/chi/v5"

	"github.com/RinaDish/currency-rates/internal/handlers"
)

type Router struct {
	l *zap.SugaredLogger
	ratesHandler handlers.RateHandler
	subscriptionHandler handlers.SubscribeHandler
}

func NewRouter(l *zap.SugaredLogger, ratesHandler handlers.RateHandler, subscriptionHandler handlers.SubscribeHandler) Router {
	return Router{
		l: l,
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