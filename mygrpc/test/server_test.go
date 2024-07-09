package test

import (
	"github.com/msprojectlb/project-common/mygrpc"
	"github.com/msprojectlb/project-common/mygrpc/registry"
	"github.com/msprojectlb/project-common/mygrpc/registry/byEtcd"
	"github.com/msprojectlb/project-common/mygrpc/test/gen"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
)

var etcdServer *clientv3.Client

func init() {
	var err error
	etcdServer, err = clientv3.New(clientv3.Config{Endpoints: []string{
		"127.0.0.1:2379",
	}})
	if err != nil {
		panic(err)
	}
}
func TestGrpcServer(t *testing.T) {

	t.Run("8080", func(t *testing.T) {
		register, err := byEtcd.NewRegister(etcdServer, 30)
		require.NoError(t, err)
		server := mygrpc.NewGrpcServer(registry.ServiceInstance{
			Name:    "appserver",
			Address: "0.0.0.0:8080",
			Weight:  40,
		}, mygrpc.WithRegistry(register))
		gen.RegisterAppServiceServer(server, &AppServer{})
		err = server.Start(":8080")
		require.NoError(t, err)
	})
	t.Run("8081", func(t *testing.T) {
		register, err := byEtcd.NewRegister(etcdServer, 30)
		require.NoError(t, err)
		server := mygrpc.NewGrpcServer(registry.ServiceInstance{
			Name:    "appserver",
			Address: "0.0.0.0:8081",
			Weight:  40,
		}, mygrpc.WithRegistry(register))
		gen.RegisterAppServiceServer(server, &AppServer{})
		err = server.Start(":8081")
		require.NoError(t, err)
	})
	t.Run("8082", func(t *testing.T) {
		register, err := byEtcd.NewRegister(etcdServer, 30)
		require.NoError(t, err)
		server := mygrpc.NewGrpcServer(registry.ServiceInstance{
			Name:    "appserver",
			Address: "0.0.0.0:8082",
			Weight:  40,
		}, mygrpc.WithRegistry(register))
		gen.RegisterAppServiceServer(server, &AppServer{})
		err = server.Start(":8082")
		require.NoError(t, err)
	})
	t.Run("8083", func(t *testing.T) {
		register, err := byEtcd.NewRegister(etcdServer, 30)
		require.NoError(t, err)
		server := mygrpc.NewGrpcServer(registry.ServiceInstance{
			Name:    "appserver",
			Address: "0.0.0.0:8083",
			Weight:  40,
		}, mygrpc.WithRegistry(register))
		gen.RegisterAppServiceServer(server, &AppServer{})
		err = server.Start(":8083")
		require.NoError(t, err)
	})
}
