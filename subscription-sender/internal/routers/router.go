package routers

import (
	"github.com/RinaDish/subscription-sender/tools"
	"github.com/go-chi/chi/v5"

	"github.com/RinaDish/subscription-sender/internal/handlers"
)

type Router struct {
	logger tools.Logger
	notifyHandler handlers.NotifyHandler
}

func NewRouter(logger tools.Logger, notifyHandler handlers.NotifyHandler) Router {
	return Router{
		logger: logger,
		notifyHandler: notifyHandler,
	}
}

func (router Router) GetRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/notify", router.notifyHandler.NotifySubscribers)

	return r
}