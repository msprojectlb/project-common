package test

import (
	"github.com/msprojectlb/project-common/config"
	myGrpc "github.com/msprojectlb/project-common/grpc"
	"github.com/msprojectlb/project-common/grpc/registry"
	"github.com/msprojectlb/project-common/grpc/registry/byEtcd"
	"github.com/msprojectlb/project-common/grpc/test/proto"
	"github.com/msprojectlb/project-common/logs"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

func Init() *clientv3.Client {
	etcdServer, err := clientv3.New(clientv3.Config{Endpoints: []string{
		"127.0.0.1:2379",
	}, DialTimeout: time.Second * 5})
	if err != nil {
		panic(err)
	}
	viper := config.NewViper(&config.ViperConf{
		AutomaticEnv: true,
		EnvPrefix:    "APP",
		ConfigPath:   []string{"./"},
		ConfigType:   "yaml",
		ConfName:     "config",
	})
	logs.InitHelper(logs.NewZapLogger(viper, logs.NewZapWriter(viper)))
	return etcdServer
}
func TestGrpcServer(t *testing.T) {
	etcdServer := Init()
	t.Run("8080", func(t *testing.T) {
		register, err := byEtcd.NewRegister(etcdServer, 30)
		require.NoError(t, err)
		server := myGrpc.NewServer(registry.ServiceInstance{
			Name:    "appserver",
			Address: "0.0.0.0:8080",
			Weight:  40,
		}, myGrpc.WithRegistry(register), myGrpc.WithLogger(logs.Helper))
		proto.RegisterTestServiceServer(server, &AppServer{})
		server.Start()
	})
	t.Run("8084", func(t *testing.T) {
		register, err := byEtcd.NewRegister(etcdServer, 30)
		require.NoError(t, err)
		server := myGrpc.NewServer(registry.ServiceInstance{
			Name:    "appserver",
			Address: "0.0.0.0:8084",
			Weight:  40,
		}, myGrpc.WithRegistry(register), myGrpc.WithLogger(logs.Helper))
		proto.RegisterTestServiceServer(server, &AppServer{})
		server.Start()
	})
	t.Run("8082", func(t *testing.T) {
		register, err := byEtcd.NewRegister(etcdServer, 30)
		require.NoError(t, err)
		server := myGrpc.NewServer(registry.ServiceInstance{
			Name:    "appserver",
			Address: "0.0.0.0:8082",
			Weight:  40,
		}, myGrpc.WithRegistry(register), myGrpc.WithLogger(logs.Helper))
		proto.RegisterTestServiceServer(server, &AppServer{})
		server.Start()
	})
	t.Run("8083", func(t *testing.T) {
		register, err := byEtcd.NewRegister(etcdServer, 30)
		require.NoError(t, err)
		server := myGrpc.NewServer(registry.ServiceInstance{
			Name:    "appserver",
			Address: "0.0.0.0:8083",
			Weight:  40,
		}, myGrpc.WithRegistry(register), myGrpc.WithLogger(logs.Helper))
		proto.RegisterTestServiceServer(server, &AppServer{})
		server.Start()
	})
}
