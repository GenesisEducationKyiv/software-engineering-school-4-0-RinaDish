package scheduler

import (
	"context"

	"github.com/RinaDish/currency-rates/internal/services"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

type Cron struct {
	l   *zap.SugaredLogger
	c 	gocron.Scheduler
}

func NewCron(l *zap.SugaredLogger, ctx context.Context, subscriptionService services.SubscriptionService) Cron {
	c, _ := gocron.NewScheduler()
	
	_, _ = c.NewJob(
		gocron.CronJob(
			"0 2 * * *",
			false,
		),
		gocron.NewTask(
			subscriptionService.NotifySubscribers, ctx,
		),
	)
	
	return Cron{
		l: l,
		c: c,
	};
}

func (cron *Cron) StartCron() (gocron.Scheduler) {
    cron.c.Start()
	cron.l.Info("Cron start")

    return cron.c
}