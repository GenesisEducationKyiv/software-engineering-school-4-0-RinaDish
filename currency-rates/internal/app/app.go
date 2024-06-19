package app

import (
	"context"
	"net/http"

	"github.com/RinaDish/currency-rates/internal/clients"
	"github.com/RinaDish/currency-rates/internal/handlers"
	"github.com/RinaDish/currency-rates/internal/repo"
	"github.com/RinaDish/currency-rates/internal/routers"
	"github.com/RinaDish/currency-rates/internal/scheduler"
	"github.com/RinaDish/currency-rates/internal/services"
	"github.com/RinaDish/currency-rates/tools"
	"github.com/go-co-op/gocron/v2"
)

type App struct {
	cfg Config
	logger   logger.Logger
	subscriptionHandler handlers.SubscribeHandler
	ratesHandler handlers.RateHandler
	subscriptionService services.SubscriptionService
	subscriptionCron gocron.Scheduler
}

func NewApp(c Config, logger logger.Logger, ctx context.Context) (*App, error) {
	nbuClient := clients.NewNBUClient(logger)
	privatClient := clients.NewPrivatClient(logger)
	rateService := services.NewRate(logger, []services.RateClient{nbuClient, privatClient})

	adminRepository, err := repo.NewAdminRepository(c.DBUrl, logger)
	if err != nil {
		return nil, err
	}

	emailSender, err := services.NewEmail(c.EmailAddress, c.EmailPass, logger)
	if err != nil {
		return nil, err
	}

	subscriptionService := services.NewSubscriptionService(logger, adminRepository, emailSender, rateService)
	ratesHandler := handlers.NewRateHandler(logger, rateService)
	subscriptionHandler := handlers.NewSubscribeHandler(logger, adminRepository)

	cron := scheduler.NewCron(logger, ctx, subscriptionService)

	subscriptionCron := cron.StartCron()

	return &App{
		cfg: c,
		logger: logger,
		subscriptionHandler: subscriptionHandler,
		ratesHandler: ratesHandler,
		subscriptionService: subscriptionService,
		subscriptionCron: subscriptionCron,
	}, nil
}

func (app *App) Run() error {
	defer func() { _ = app.subscriptionCron.Shutdown() }()

	r := routers.NewRouter(app.logger, app.ratesHandler, app.subscriptionHandler)

	app.logger.Info("app run")
	
	return http.ListenAndServe(app.cfg.Address, r.GetRouter())
}
