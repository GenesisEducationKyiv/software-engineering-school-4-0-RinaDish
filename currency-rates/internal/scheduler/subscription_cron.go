package scheduler

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/RinaDish/currency-rates/tools"
)

type Cron struct {
	logger  tools.Logger
	scheduler  gocron.Scheduler
}

func NewCron(logger tools.Logger) Cron {
	scheduler, _ := gocron.NewScheduler()
	
	return Cron{
		logger: logger,
		scheduler: scheduler,
	};
}

func (cron *Cron) RegisterTask(schedule string, task gocron.Task) {
	_, _ = cron.scheduler.NewJob(
		gocron.CronJob(
			schedule,
			false,
		),
		task,
	)
}

func (cron *Cron) StartCron() (gocron.Scheduler) {
    cron.scheduler.Start()
	cron.logger.Info("Cron start")

    return cron.scheduler
}