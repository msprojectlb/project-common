package handler

import (
	"github.com/google/wire"
	myGrpc "github.com/msprojectlb/project-common/grpc"
	"github.com/msprojectlb/project-common/grpc/registry"
	"github.com/msprojectlb/project-common/grpc/registry/byEtcd"
	"github.com/msprojectlb/project-common/logs"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"io"
	"time"
)

// LoggerSet 全局日志
var LoggerSet = wire.NewSet(logs.NewZapLogger, logs.NewZapWriter, wire.Bind(new(io.Writer), new(*logs.ZapWriter)))

// RegisterByEtcdSet etcd注册中心
var RegisterByEtcdSet = wire.NewSet(byEtcd.NewRegister, NewEtcdClient, wire.Value(30))

// GrpcServerSet grpc服务端
var GrpcServerSet = wire.NewSet(myGrpc.NewServer, registry.NewServiceInstance, GetGrpcServerOptions, RegisterByEtcdSet)

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
