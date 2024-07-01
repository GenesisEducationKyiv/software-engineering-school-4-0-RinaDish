package app

import (
	"context"
	"net/http"

	"github.com/RinaDish/subscription-sender/internal/handlers"
	"github.com/RinaDish/subscription-sender/internal/routers"
	"github.com/RinaDish/subscription-sender/internal/services"
	"github.com/RinaDish/subscription-sender/tools"
)

type App struct {
	cfg Config
	logger   tools.Logger
	router routers.Router
}

func NewApp(c Config, logger tools.Logger, ctx context.Context) (*App, error) {
	emailSender, err := services.NewEmail(c.EmailAddress, c.EmailPass, logger)
	if err != nil {
		return nil, err
	}

	subscriptionService := services.NewSubscriptionService(logger, emailSender)

	notifyHandler := handlers.NewNotifyHandler(logger, subscriptionService)

	router := routers.NewRouter(logger, notifyHandler)

	return &App{
		cfg: c,
		logger: logger,
		router: router,
	}, nil
}

func (app *App) Run() error {
	app.logger.Info("app run")

	return http.ListenAndServe(app.cfg.Address, app.router.GetRouter())
}
