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

func TestGrpcServer(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{
			"127.0.0.1:2379",
		},
	})
	require.NoError(t, err)
	register, err := byEtcd.NewRegister(client, 30)
	require.NoError(t, err)
	server := mygrpc.NewGrpcServer(registry.ServiceInstance{
		Name:    "appserver",
		Address: "0.0.0.0:8080",
	}, mygrpc.WithRegistry(register))
	gen.RegisterAppServiceServer(server, &AppServer{})
	err = server.Start(":8080")
	require.NoError(t, err)
}

func TestGrpcServer2(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{
			"127.0.0.1:2379",
		},
	})
	require.NoError(t, err)
	register, err := byEtcd.NewRegister(client, 30)
	require.NoError(t, err)
	server := mygrpc.NewGrpcServer(registry.ServiceInstance{
		Name:    "appserver",
		Address: "0.0.0.0:8081",
	}, mygrpc.WithRegistry(register))
	gen.RegisterAppServiceServer(server, &AppServer{})
	err = server.Start(":8081")
	require.NoError(t, err)
}
