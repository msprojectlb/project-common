package test

import (
	"context"
	"github.com/msprojectlb/project-common/mygrpc"
	"github.com/msprojectlb/project-common/mygrpc/registry/byEtcd"
	"github.com/msprojectlb/project-common/mygrpc/test/gen"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

func TestGrpcClient(t *testing.T) {
	c, err := clientv3.New(clientv3.Config{Endpoints: []string{
		"127.0.0.1:2379",
	}})
	require.NoError(t, err)
	register, err := byEtcd.NewRegister(c, 30)
	require.NoError(t, err)
	dial, err := grpc.NewClient("etcd:///appserver", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithResolvers(mygrpc.NewGrpcResolverBuilder(register)))
	require.NoError(t, err)
	defer dial.Close()
	client := gen.NewAppServiceClient(dial)
	hello, err := client.Hello(context.Background(), &gen.HelloReq{Name: "张三"})
	require.NoError(t, err)
	t.Log(hello)
	hello2, err := client.Hello(context.Background(), &gen.HelloReq{Name: "张三2"})
	require.NoError(t, err)
	t.Log(hello2)
	hello3, err := client.Hello(context.Background(), &gen.HelloReq{Name: "张三3"})
	require.NoError(t, err)
	t.Log(hello3)
}
