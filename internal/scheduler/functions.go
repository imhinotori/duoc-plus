package scheduler

import (
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/redis/go-redis/v9"
	"time"
)

func (w *Worker) updateTokens() {
	log.Debug("Searching for tokens to update...")
	ctx := context.TODO()

	tasks, err := w.db.Scheduler.ZRangeByScoreWithScores(ctx, "tasks", &redis.ZRangeBy{
		Min: "-inf",
		Max: fmt.Sprintf("%f", float64(time.Now().Unix())),
	}).Result()
	if err != nil {
		log.Debugf("error getting tasks: %s", err)
		return
	}

	for _, task := range tasks {

		w.db.Scheduler.ZRem(ctx, "tasks", task.Member)
	}

	log.Debug("Tokens updated")
}
