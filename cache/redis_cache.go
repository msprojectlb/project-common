package cache

import (
	"context"
	"errors"
	"github.com/msprojectlb/project-common/db"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	*db.RedisDb
}

func NewRedisCache(rdb *db.RedisDb) Cache {
	return &RedisCache{RedisDb: rdb}
}

func (rc *RedisCache) Put(ctx context.Context, key, value string, expire time.Duration) error {
	err := rc.Rdb.Set(ctx, key, value, expire).Err()
	return err
}
func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := rc.Rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return result, err
}
