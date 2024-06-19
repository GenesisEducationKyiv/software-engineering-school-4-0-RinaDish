package scheduler

import (
	"context"

	"github.com/RinaDish/currency-rates/internal/services"
	"github.com/go-co-op/gocron/v2"
	"github.com/RinaDish/currency-rates/tools"
)

type Cron struct {
	logger  logger.Logger
	cron  gocron.Scheduler
}

func NewCron(logger logger.Logger, ctx context.Context, subscriptionService services.SubscriptionService) Cron {
	cron, _ := gocron.NewScheduler()
	
	_, _ = cron.NewJob(
		gocron.CronJob(
			"0 2 * * *",
			false,
		),
		gocron.NewTask(
			subscriptionService.NotifySubscribers, ctx,
		),
	)
	
	return Cron{
		logger: logger,
		cron: cron,
	};
}

func (cron *Cron) StartCron() (gocron.Scheduler) {
    cron.cron.Start()
	cron.logger.Info("Cron start")

    return cron.cron
}