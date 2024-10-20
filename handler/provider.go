package handler

import (
	"github.com/google/wire"
	"github.com/msprojectlb/project-common/config"
	myGrpc "github.com/msprojectlb/project-common/grpc"
	"github.com/msprojectlb/project-common/grpc/registry"
	"github.com/msprojectlb/project-common/grpc/registry/byEtcd"
	"github.com/msprojectlb/project-common/logs"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// Logger 全局日志
var Logger = wire.NewSet(logs.NewZapLogger, config.NewViper, logs.NewZapWriter)

// RegisterByEtcd etcd注册中心
var RegisterByEtcd = wire.NewSet(byEtcd.NewRegister, NewEtcdClient, wire.Value(30))

// GrpcServer grpc服务端
var GrpcServer = wire.NewSet(myGrpc.NewServer, registry.NewServiceInstance, GetGrpcServerOptions, RegisterByEtcd)

// NewEtcdClient etcd客户端
func NewEtcdClient(conf *viper.Viper) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   conf.GetStringSlice("etcd.addr"),
		DialTimeout: time.Second * 10,
	})
}

// GetGrpcServerOptions grpc服务配置
func GetGrpcServerOptions(r registry.Registry, l *logs.ZapLogger) []myGrpc.Options {
	return []myGrpc.Options{
		myGrpc.WithRegistry(r),
		myGrpc.WithLogger(l),
	}
}
