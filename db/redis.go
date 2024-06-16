package db

import (
	"context"
	"fmt"
	"github.com/msprojectlb/project-common/config"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type RedisDb struct {
	Rdb *redis.Client
}

func NewRedisDb(rdb *redis.Client) *RedisDb {
	return &RedisDb{Rdb: rdb}
}

// NewSingleRdb 单节点
func NewSingleRdb(conf config.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.Db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf(fmt.Sprintf("redis 连接失败: %s", err.Error()))
	}
	return rdb
}
