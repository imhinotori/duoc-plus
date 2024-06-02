package database

import (
	"context"
	"github.com/charmbracelet/log"
	"github.com/imhinotori/duoc-plus/internal/config"
	"github.com/redis/go-redis/v9"
	"time"
)

func newRedisClient(cfg *config.Config, db int) (*redis.Client, error) {
	opts, err := redis.ParseURL(cfg.Redis.ConnectionURI)
	if err != nil {
		log.Debugf("error parsing connection URI: %s", err)
		return nil, err
	}

	opts.DB = db

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	client := redis.NewClient(opts)

	_, err = client.Ping(ctx).Result()
	if err != nil {
		log.Debugf("error connecting to Users: %s", err)
		return nil, err
	}

	return client, nil
}
