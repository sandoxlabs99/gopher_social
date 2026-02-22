package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sandoxlabs99/gopher_social/internal/models"

	"github.com/redis/go-redis/v9"
)

type UserStore struct {
	rdb *redis.Client
}

const UserExpTime = time.Minute

func (rds *UserStore) Get(ctx context.Context, userID int64) (*models.User, error) {
	cacheKey := fmt.Sprintf("user-%v", userID)

	data, err := rds.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user models.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (rds *UserStore) Set(ctx context.Context, user *models.User) error {
	cacheKey := fmt.Sprintf("user-%v", user.ID)

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return rds.rdb.SetEx(ctx, cacheKey, json, UserExpTime).Err()
}
