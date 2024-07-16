package app

import (
	"context"
	"net/http"

	"github.com/RinaDish/currency-rates/internal/clients"
	"github.com/RinaDish/currency-rates/internal/handlers"
	"github.com/RinaDish/currency-rates/internal/repo"
	"github.com/RinaDish/currency-rates/internal/routers"
	"github.com/RinaDish/currency-rates/internal/services"
	"github.com/RinaDish/currency-rates/internal/workers"
	"github.com/RinaDish/currency-rates/internal/queue"
	"github.com/RinaDish/currency-rates/tools"
	"github.com/go-co-op/gocron/v2"
	"github.com/nats-io/nats.go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	cfg Config
	logger   tools.Logger
	subscriptionHandler handlers.SubscribeHandler
	ratesHandler handlers.RateHandler
	subscriptionService services.SubscriptionService
	subscriptionCron  workers.Cron
	router routers.Router
	db *gorm.DB
	queue *queue.Queue
}

func NewApp(cfg Config, logger tools.Logger, ctx context.Context) (*App, error) {
	nbuClient := clients.NewNBUClient(logger)
	privatClient := clients.NewPrivatClient(logger)

	nbuChain := clients.NewBaseChain(nbuClient)
	privatChain := clients.NewBaseChain(privatClient)
	nbuChain.SetNext(privatChain)
	
	rateService := services.NewRate(logger, nbuChain)

	db, err := gorm.Open(postgres.Open(cfg.DBURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	adminRepository := repo.NewAdminRepository(db, logger)

	nats, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		return nil, err
	}

	natsbroker := queue.NewNATSBroker(nats)

	queue := queue.NewQueue(natsbroker, cfg.SubscriptionTopicName, logger)

	subscriptionService := services.NewSubscriptionService(logger, adminRepository, queue, rateService)
	ratesHandler := handlers.NewRateHandler(logger, rateService)
	subscriptionHandler := handlers.NewSubscribeHandler(logger, adminRepository)

	cron := workers.NewCron(logger)
	task := gocron.NewTask(subscriptionService.NotifySubscribers, ctx)
	
	cron.RegisterTask("0 2 * * *", task)

	router := routers.NewRouter(logger, ratesHandler, subscriptionHandler)

	return &App{
		cfg: cfg,
		logger: logger,
		subscriptionHandler: subscriptionHandler,
		ratesHandler: ratesHandler,
		subscriptionService: subscriptionService,
		subscriptionCron: cron,
		router: router,
		db: db,
		queue: queue,
	}, nil
}

func (app *App) Run() error {
	subscriptionCron := app.subscriptionCron.StartCron()

	defer func() { 
			_ = subscriptionCron.Shutdown() 
			
			if db, err := app.db.DB(); err == nil {
			_ = db.Close()
			}

			_ = app.queue.Broker.Drain()
		}()

	app.logger.Info("app run")
	
	return http.ListenAndServe(app.cfg.Address, app.router.GetRouter())
}
