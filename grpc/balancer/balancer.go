package balancer

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

func init() {
	//注册加权轮询均衡器
	balancer.Register(base.NewBalancerBuilder(WeightPollingBalancerName, &WeightPollingBalancer{}, base.Config{HealthCheck: true}))
	//注册轮询均衡器
	balancer.Register(base.NewBalancerBuilder(PollingBalancerName, &PollingBalancer{}, base.Config{HealthCheck: true}))
}
