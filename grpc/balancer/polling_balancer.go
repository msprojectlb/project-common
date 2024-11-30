package balancer

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math"
	"sync/atomic"
)

func init() {
	balancer.Register(base.NewBalancerBuilder(PollingBalancerName, &PollingBalancer{}, base.Config{HealthCheck: true}))
}

const PollingBalancerName = "polling-balancer"

// PollingBalancer 轮询
type PollingBalancer struct {
}

func (p *PollingBalancer) Build(info base.PickerBuildInfo) balancer.Picker {
	var res PollingPicker
	res.connection = make([]balancer.SubConn, 0, len(info.ReadySCs))

	for conn := range info.ReadySCs {
		res.connection = append(res.connection, conn)
	}

	res.length = int32(len(res.connection))
	res.index = -1
	return &res
}

type PollingPicker struct {
	index      int32
	length     int32
	connection []balancer.SubConn
}

func (p *PollingPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	var res balancer.PickResult
	if p.length == 0 {
		return res, balancer.ErrNoSubConnAvailable
	}
	idx := atomic.AddInt32(&p.index, 1)
	if idx > math.MaxInt32 {
		atomic.CompareAndSwapInt32(&p.index, idx, -1)
		idx = atomic.AddInt32(&p.index, 1)
	}
	res.SubConn = p.connection[idx%p.length]
	return res, nil
}
