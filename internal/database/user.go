package database

import (
	"context"
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/imhinotori/duoc-plus/internal/common"
	"time"
)

func (db *Database) GetUserFromSessionId(id string) (*common.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	data, err := db.Users.Get(ctx, id).Result()
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
