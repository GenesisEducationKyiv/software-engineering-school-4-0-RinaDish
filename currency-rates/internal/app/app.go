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
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

type App struct {
	cfg Config
	l   *zap.SugaredLogger
	subscriptionHandler handlers.SubscribeHandler
	ratesHandler handlers.RateHandler
	subscriptionService services.SubscriptionService
	subscriptionCron gocron.Scheduler
}

func NewApp(c Config, l *zap.SugaredLogger, ctx context.Context) (*App, error) {
	nbuClient := clients.NewNBUClient(l)
	privatClient := clients.NewPrivatClient(l)
	rateService := services.NewRate(l, []services.RateClient{nbuClient, privatClient})

	adminRepository, err := repo.NewAdminRepository(c.DBUrl, l)
	if err != nil {
		return nil, err
	}

	emailSender, err := services.NewEmail(c.EmailAddress, c.EmailPass, l)
	if err != nil {
		return nil, err
	}

	subscriptionService := services.NewSubscriptionService(l, adminRepository, emailSender, rateService)
	ratesHandler := handlers.NewRateHandler(l, rateService)
	subscriptionHandler := handlers.NewSubscribeHandler(l, adminRepository)

	cron := scheduler.NewCron(l, ctx, subscriptionService)

	subscriptionCron := cron.StartCron()

	return &App{
		cfg: c,
		l:   l,
		subscriptionHandler: subscriptionHandler,
		ratesHandler: ratesHandler,
		subscriptionService: subscriptionService,
		subscriptionCron: subscriptionCron,
	}, nil
}

func (app *App) Run() error {
	defer func() { _ = app.subscriptionCron.Shutdown() }()

	r := routers.NewRouter(app.l, app.ratesHandler, app.subscriptionHandler)

	app.l.Info("app run")
	
	return http.ListenAndServe(app.cfg.Address, r.GetRouter())
}
