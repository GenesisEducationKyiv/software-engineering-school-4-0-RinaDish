package app

import (
	"context"
	"net/http"

	"github.com/RinaDish/subscription-sender/internal/handlers"
	"github.com/RinaDish/subscription-sender/internal/queue"
	"github.com/RinaDish/subscription-sender/internal/repo"
	"github.com/RinaDish/subscription-sender/internal/routers"
	"github.com/RinaDish/subscription-sender/internal/services"
	"github.com/RinaDish/subscription-sender/tools"
	"github.com/RinaDish/subscription-sender/workers"
	"github.com/go-co-op/gocron/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/nats-io/nats.go"
)

type App struct {
	cfg              Config
	logger           tools.Logger
	subscriptionCron workers.Cron
	router           routers.Router
	queue            *queue.Queue
	db               *gorm.DB
	ctx              context.Context
}

func NewApp(cfg Config, logger tools.Logger, ctx context.Context) (*App, error) {
	emailSender, err := services.NewEmail(cfg.EmailAddress, cfg.EmailPass, logger)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(cfg.DBURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	repo := repo.NewAdminRepository(db, logger)

	notificationService := services.NewNotificationService(logger, emailSender, repo)

	messagesService := services.NewMessagesService(repo, logger)

	healthCheckHandler := handlers.NewHealthCheckHandler(logger)

	router := routers.NewRouter(logger, healthCheckHandler)

	nats, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		return nil, err
	}

	queue := queue.NewQueue(nats, cfg.SubscriptionTopicName, messagesService, logger)

	cron := workers.NewCron(logger)
	task := gocron.NewTask(notificationService.NotifySubscribers, ctx)

	cron.RegisterTask("0 * * * *", task)

	return &App{
		cfg:              cfg,
		logger:           logger,
		router:           router,
		subscriptionCron: cron,
		queue:            queue,
		db:               db,
		ctx:              ctx,
	}, nil
}

func (app *App) Run() error {
	subscriptionCron := app.subscriptionCron.StartCron()

	defer func() {
		_ = subscriptionCron.Shutdown()

		if db, err := app.db.DB(); err == nil {
			_ = db.Close()
		}

		_ = app.queue.QueueConn.Drain()
	}()

	err := app.queue.Subscribe(app.ctx)

	if err != nil {
		app.logger.Error("Queue subscribe method faild")
	}

	if err != nil {
		app.logger.Errorf("err")
	}

	app.logger.Info("app run")

	return http.ListenAndServe(app.cfg.Address, app.router.GetRouter())
}
