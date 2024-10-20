package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"time"
)

type RedisDb struct {
	Rdb *redis.Client
}

func NewRedisDb(rdb *redis.Client) *RedisDb {
	return &RedisDb{Rdb: rdb}
}

// NewSingleRdb 单节点
func NewSingleRdb(viper *viper.Viper) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(err)
	}
	return rdb
}
