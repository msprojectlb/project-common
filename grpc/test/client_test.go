package test

import (
	"context"
	"fmt"
	myGrpc "github.com/msprojectlb/project-common/grpc"
	mybalancer "github.com/msprojectlb/project-common/grpc/balancer"
	"github.com/msprojectlb/project-common/grpc/registry/byEtcd"
	"github.com/msprojectlb/project-common/grpc/test/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
	"sync"
	"testing"
)

func TestGrpcClientWithPollingBalance(t *testing.T) {
	etcdClient := Init()
	register, err := byEtcd.NewRegister(etcdClient, 30)
	require.NoError(t, err)
	//初始化负载均衡器
	balancer.Register(base.NewBalancerBuilder(mybalancer.PollingBalancerName, &mybalancer.PollingBalancer{}, base.Config{HealthCheck: true}))

	dial, err := grpc.NewClient(
		"etcd:///appserver",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(myGrpc.NewResolverBuilder(register)),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy":"%s"}`, mybalancer.PollingBalancerName)))
	require.NoError(t, err)
	defer dial.Close()
	client := proto.NewTestServiceClient(dial)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			hello, err := client.Hello(context.Background(), &proto.HelloReq{Name: "张三__" + strconv.Itoa(i)})
			require.NoError(t, err)
			t.Log(hello)
		}(i)
	}
	wg.Wait()
}

func TestGrpcClientWithWeightPollingBalance(t *testing.T) {
	etcdClient := Init()
	register, err := byEtcd.NewRegister(etcdClient, 30)
	require.NoError(t, err)
	//初始化加权轮询均衡器
	balancer.Register(base.NewBalancerBuilder(mybalancer.WeightPollingBalancerName, &mybalancer.WeightPollingBalancer{}, base.Config{HealthCheck: true}))

	dial, err := grpc.NewClient(
		"etcd:///appserver",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(myGrpc.NewResolverBuilder(register)),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy":"%s"}`, mybalancer.WeightPollingBalancerName)))
	require.NoError(t, err)
	defer dial.Close()
	client := proto.NewTestServiceClient(dial)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			hello, err := client.Hello(context.Background(), &proto.HelloReq{Name: "张三__" + strconv.Itoa(i)})
			require.NoError(t, err)
			t.Log(hello)
		}(i)
	}
	wg.Wait()
}
