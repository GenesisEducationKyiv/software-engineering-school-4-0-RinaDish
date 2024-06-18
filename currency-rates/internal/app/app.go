package app

import (
	"context"
	"net/http"

	"github.com/RinaDish/currency-rates/internal/clients"
	"github.com/RinaDish/currency-rates/internal/handlers"
	"github.com/RinaDish/currency-rates/internal/repo"
	"github.com/RinaDish/currency-rates/internal/routers"
	"github.com/RinaDish/currency-rates/internal/services"
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

func (app *App) initCron(ctx context.Context, subscriptionService services.SubscriptionService) gocron.Scheduler {
	s, _ := gocron.NewScheduler()
	
	_, _ = s.NewJob(
		gocron.CronJob(
			"0 2 * * *",
			false,
		),
		gocron.NewTask(
			subscriptionService.NotifySubscribers, ctx,
		),
	)
	
	return s;
}

func (app *App) startCron(ctx context.Context, subscriptionService services.SubscriptionService) (gocron.Scheduler) {
    s := app.initCron(ctx, subscriptionService)
    s.Start()
	app.l.Info("Cron start")

    return s
}

func (app *App) Run(ctx context.Context) error {
	nbuClient := clients.NewNBUClient(app.l)
	privatClient := clients.NewPrivatClient(app.l)
	rateService := services.NewRate(app.l, []services.RateClient{nbuClient, privatClient})

	adminRepository, err := repo.NewAdminRepository(app.cfg.DBUrl, app.l)
	if err != nil {
		return err
	}

	emailSender, err := services.NewEmail(app.cfg.EmailAddress, app.cfg.EmailPass, app.l)
	if err != nil {
		return err
	}

	subscriptionService := services.NewSubscriptionService(app.l, adminRepository, emailSender, rateService)
	ratesHandler := handlers.NewRateHandler(app.l, rateService)
	subscriptionHandler := handlers.NewSubscribeHandler(app.l, adminRepository)

	s := app.startCron(ctx, subscriptionService)

	defer func() { _ = s.Shutdown() }()

	r := routers.NewRouter(app.l, ratesHandler, subscriptionHandler)

	app.l.Info("app run")
	
	return http.ListenAndServe(app.cfg.Address, r.GetRouter())
}
