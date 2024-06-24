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
	logger   tools.Logger
	subscriptionHandler handlers.SubscribeHandler
	ratesHandler handlers.RateHandler
	subscriptionService services.SubscriptionService
	subscriptionCron  scheduler.Cron
	router routers.Router
}

func NewApp(c Config, logger tools.Logger, ctx context.Context) (*App, error) {
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

	cron := scheduler.NewCron(logger)
	task := gocron.NewTask(subscriptionService.NotifySubscribers, ctx)
	
	cron.RegisterTask("0 2 * * *", task)

	router := routers.NewRouter(logger, ratesHandler, subscriptionHandler)

	return &App{
		cfg: c,
		logger: logger,
		subscriptionHandler: subscriptionHandler,
		ratesHandler: ratesHandler,
		subscriptionService: subscriptionService,
		subscriptionCron: cron,
		router: router,
	}, nil
}

func (app *App) Run() error {
	subscriptionCron := app.subscriptionCron.StartCron()

	defer func() { _ = subscriptionCron.Shutdown() }()

	app.logger.Info("app run")
	
	return http.ListenAndServe(app.cfg.Address, app.router.GetRouter())
}
