package app

import (
	"context"
	"net/http"

	"github.com/RinaDish/currency-rates/internal/clients"
	"github.com/RinaDish/currency-rates/internal/handlers"
	"github.com/RinaDish/currency-rates/internal/repo"
	"github.com/RinaDish/currency-rates/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

type App struct {
	cfg Config
	l   *zap.SugaredLogger
}

func NewApp(c Config, l *zap.SugaredLogger) *App {
	return &App{
		cfg: c,
		l:   l,
	}
}

func (app *App) Run(ctx context.Context) error {
	nbuClient := clients.NewNBUClient(app.l)
	privatClient := clients.NewPrivatClient(app.l)
	rateService := services.NewRate(app.l, []services.RateClient{nbuClient, privatClient})

	adminRepository, err := repo.NewAdminRepository(app.cfg.DBUrl, app.l)
	if err != nil {
		return err
	}

	subscriptionSender := services.NewEmail(app.cfg.EmailAddress, app.cfg.EmailPass, app.l)
	subscriptionService := services.NewSubscriptionService(app.l, adminRepository, subscriptionSender, rateService)
	ratesHandler := handlers.NewRateHandler(app.l, rateService)
	subscriptionHandler := handlers.NewSubscribeHandler(app.l, adminRepository)

	s, _ := gocron.NewScheduler()
	defer func() { _ = s.Shutdown() }()
	
	_, _ = s.NewJob(
		gocron.CronJob(
			"0 2 * * *",
			false,
		),
		gocron.NewTask(
			subscriptionService.NotifySubscribers, ctx,
		),
	)

	s.Start()

	r := chi.NewRouter()
	r.Get("/rate", ratesHandler.GetCurrentRate)
	r.Post("/subscribe", subscriptionHandler.CreateSubscription)

	app.l.Info("app run")
	return http.ListenAndServe(app.cfg.Address, r)
}
