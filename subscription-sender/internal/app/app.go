package app

import (
	"context"
	"net/http"

	"github.com/RinaDish/subscription-sender/internal/handlers"
	"github.com/RinaDish/subscription-sender/internal/queue/subscription_notifier"
	"github.com/RinaDish/subscription-sender/internal/queue"
	"github.com/RinaDish/subscription-sender/internal/routers"
	"github.com/RinaDish/subscription-sender/internal/services"
	"github.com/RinaDish/subscription-sender/tools"

	"github.com/nats-io/nats.go"
)

type App struct {
	cfg Config
	logger   tools.Logger
	router routers.Router
	queue *subscriptionnotifier.SubscriptionNotifierConsumer
	ctx context.Context
}

func NewApp(cfg Config, logger tools.Logger, ctx context.Context) (*App, error) {
	emailSender, err := services.NewEmail(cfg.EmailAddress, cfg.EmailPass, logger)
	if err != nil {
		return nil, err
	}

	subscriptionService := services.NewSubscriptionService(logger, emailSender)

	healthCheckHandler := handlers.NewHealthCheckHandler(logger)

	router := routers.NewRouter(logger, healthCheckHandler)

	nats, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		return nil, err
	}

	natsbroker := queue.NewNATSBroker(nats)

	queue := subscriptionnotifier.NewSubscriptionNotifierConsumer(natsbroker, cfg.SubscriptionTopicName, subscriptionService, logger)

	return &App{
		cfg: cfg,
		logger: logger,
		router: router,
		queue: queue,
		ctx: ctx,
	}, nil
}

func (app *App) Run() error {
	app.logger.Info("app run")
	
	defer func() { 
		_ = app.queue.Broker.Drain()
	}()

	if err := app.queue.ConsumeSubscriptionEvent(app.ctx); err != nil {
        app.logger.Error("Queue subscribe method failed")
		return err
    }
	
	return http.ListenAndServe(app.cfg.Address, app.router.GetRouter())
}
