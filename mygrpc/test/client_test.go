package test

import (
	"context"
	"github.com/msprojectlb/project-common/mygrpc"
	mybalancer "github.com/msprojectlb/project-common/mygrpc/balancer"
	"github.com/msprojectlb/project-common/mygrpc/registry/byEtcd"
	"github.com/msprojectlb/project-common/mygrpc/test/gen"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
	"sync"
	"testing"
)

var etcdClient *clientv3.Client

func init() {
	var err error
	etcdClient, err = clientv3.New(clientv3.Config{Endpoints: []string{
		"127.0.0.1:2379",
	}})
	if err != nil {
		panic(err)
	}
}

func TestGrpcClient(t *testing.T) {
	register, err := byEtcd.NewRegister(etcdClient, 30)
	require.NoError(t, err)
	//初始化负载均衡器
	balancer.Register(base.NewBalancerBuilder("LOADBALANCING", &mybalancer.PollingBalancer{}, base.Config{HealthCheck: true}))

	dial, err := grpc.NewClient("etcd:///appserver", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithResolvers(mygrpc.NewGrpcResolverBuilder(register)), grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy":"LOADBALANCING"}`))
	require.NoError(t, err)
	defer dial.Close()
	client := gen.NewAppServiceClient(dial)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			hello, err := client.Hello(context.Background(), &gen.HelloReq{Name: "张三__" + strconv.Itoa(i)})
			require.NoError(t, err)
			t.Log(hello)
		}(i)
	}
	wg.Wait()
}

func TestGrpcClientWithWeightPollingBalance(t *testing.T) {
	register, err := byEtcd.NewRegister(etcdClient, 30)
	require.NoError(t, err)
	//初始化加权轮询均衡器
	balancer.Register(base.NewBalancerBuilder("weight-polling-balance", &mybalancer.WeightPollingBalancer{}, base.Config{HealthCheck: true}))

	dial, err := grpc.NewClient("etcd:///appserver", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithResolvers(mygrpc.NewGrpcResolverBuilder(register)), grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy":"weight-polling-balance"}`))
	require.NoError(t, err)
	defer dial.Close()
	client := gen.NewAppServiceClient(dial)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			hello, err := client.Hello(context.Background(), &gen.HelloReq{Name: "张三__" + strconv.Itoa(i)})
			require.NoError(t, err)
			t.Log(hello)
		}(i)
	}
	wg.Wait()
}
