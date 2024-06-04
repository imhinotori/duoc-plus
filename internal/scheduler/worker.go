package scheduler

import (
	"github.com/charmbracelet/log"
	"github.com/go-co-op/gocron/v2"
	"github.com/imhinotori/duoc-plus/internal/database"
)

const (
	jobFrequency = "* * * * *"
)

type Worker struct {
	db        *database.Database
	scheduler gocron.Scheduler
	job       gocron.Job
}

func New(db *database.Database) *Worker {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil
	}

	return &Worker{
		db:        db,
		scheduler: scheduler,
	}
}

func (w *Worker) Start() error {
	job, err := w.scheduler.NewJob(
		gocron.CronJob(jobFrequency, false),
		gocron.NewTask(w.Process),
	)

	if err != nil {
		return err
	}

	w.scheduler.Start()
	w.job = job
	return nil
}

func (w *Worker) Process() {
	log.Debug("Processing...")
	w.updateTokens()
	log.Debug("Finished processing")
}
