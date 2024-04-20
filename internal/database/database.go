package database

import (
	"github.com/charmbracelet/log"
	"github.com/imhinotori/duoc-plus/internal/config"
	"github.com/redis/go-redis/v9"
)

type Database struct {
	Redis redis.Client
}

func New(cfg *config.Config) (*Database, error) {
	redisClient, err := newRedisClient(cfg)
	if err != nil {
		log.Debugf("error creating Redis client: %s", err)
		return nil, err
	}

	return &Database{
		Redis: *redisClient,
	}, nil

}
