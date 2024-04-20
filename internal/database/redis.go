package database

import (
	"context"
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/imhinotori/duoc-plus/internal/common"
	"github.com/imhinotori/duoc-plus/internal/config"
	"github.com/redis/go-redis/v9"
	"time"
)

func newRedisClient(cfg *config.Config) (*redis.Client, error) {
	opts, err := redis.ParseURL(cfg.Redis.ConnectionURI)
	if err != nil {
		log.Debugf("error parsing connection URI: %s", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	client := redis.NewClient(opts)

	_, err = client.Ping(ctx).Result()
	if err != nil {
		log.Debugf("error connecting to Redis: %s", err)
		return nil, err
	}

	return client, nil
}

func (db *Database) GetUserFromSessionId(id string) (*common.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	data, err := db.Redis.Get(ctx, id).Result()
	if err != nil {
		log.Debugf("error getting user from session ID: %s", err)
		return nil, err
	}

	var usr common.User

	err = json.Unmarshal([]byte(data), &usr)
	if err != nil {
		log.Debugf("error unmarshalling user data: %s", err)
		return nil, err
	}

	log.Debugf("user data: %v", usr)

	return &usr, nil
}
