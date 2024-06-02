package database

import (
	"github.com/charmbracelet/log"
	"github.com/imhinotori/duoc-plus/internal/config"
	"github.com/redis/go-redis/v9"
)

type Database struct {
	Users     redis.Client
	Scheduler redis.Client
}

func New(cfg *config.Config) (*Database, error) {
	redisClient, err := newRedisClient(cfg, 0)
	if err != nil {
		log.Debugf("error creating Users client: %s", err)
		return nil, err
	}

	schedulerRedisClient, err := newRedisClient(cfg, 5)
	if err != nil {
		log.Debugf("error creating Users client: %s", err)
		return nil, err
	}

	return &Database{
		Users:     *redisClient,
		Scheduler: *schedulerRedisClient,
	}, nil

}
