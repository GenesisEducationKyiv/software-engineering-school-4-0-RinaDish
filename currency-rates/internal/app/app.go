package app

import (
	"context"
	"net/http"

	"github.com/RinaDish/currency-rates/internal/clients"
	"github.com/RinaDish/currency-rates/internal/handlers"
	"github.com/RinaDish/currency-rates/internal/repo"
	"github.com/RinaDish/currency-rates/internal/routers"
	"github.com/RinaDish/currency-rates/internal/services"
	"github.com/RinaDish/currency-rates/internal/scheduler"
	"go.uber.org/zap"
)

type App struct {
	cfg Config
	l   *zap.SugaredLogger
	subscriptionHandler handlers.SubscribeHandler
	ratesHandler handlers.RateHandler
	subscriptionService services.SubscriptionService
}

func NewApp(c Config, l *zap.SugaredLogger) (*App, error) {
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

	return &App{
		cfg: c,
		l:   l,
		subscriptionHandler: subscriptionHandler,
		ratesHandler: ratesHandler,
		subscriptionService: subscriptionService,
	}, nil
}

func (app *App) Run(ctx context.Context) error {
	cron := scheduler.NewCron(app.l, ctx, app.subscriptionService)

	s := cron.StartCron()

	defer func() { _ = s.Shutdown() }()

	r := routers.NewRouter(app.l, app.ratesHandler, app.subscriptionHandler)

	app.l.Info("app run")
	
	return http.ListenAndServe(app.cfg.Address, r.GetRouter())
}
