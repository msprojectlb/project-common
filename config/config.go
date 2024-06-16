package config

import (
	"google.golang.org/grpc"
	"gorm.io/gorm/logger"
	"time"
)

type GRPCConfig struct {
	Addr         string
	RegisterFunc func(*grpc.Server)
}
type EtcdConfig struct {
	Addrs []string
}
type MysqlConfig struct {
	User         string
	Pwd          string
	Ip           string
	Db           string
	CharSet      string
	Port         int
	MaxIdleConns int           //空闲连接池最大数量
	MaxOpenConns int           //最大打开的连接数
	MaxIdleTime  time.Duration //最大空闲连接时间
	MaxLifetime  time.Duration //连接可复用的最长时间
}
type GormConfig struct {
	IgnoreRecordNotFoundError bool //Ignore ErrRecordNotFound error for logger
	ParameterizedQueries      bool //Don't include params in the SQL log
	Colorful                  bool
	SlowThreshold             time.Duration   //Slow SQL threshold
	LogLevel                  logger.LogLevel //Log level
}
type RedisConfig struct {
	Password string
	Addr     string
	Db       int
}

type JWTConfig struct {
	AccessExp     time.Duration
	RefreshExp    time.Duration
	AccessSecret  string
	RefreshSecret string
}
