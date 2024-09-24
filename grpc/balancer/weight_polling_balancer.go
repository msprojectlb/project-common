package balancer

import (
	"cmp"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"sync"
)

const WeightPollingBalancerName = "weight-polling-balancer"

type node struct {
	conn            balancer.SubConn
	weight          int
	currentWeight   int
	effectiveWeight int
}

// WeightPollingBalancer 加权平滑轮询
type WeightPollingBalancer struct {
}

func (w *WeightPollingBalancer) Build(info base.PickerBuildInfo) balancer.Picker {
	var res WeightPollingPicker
	res.nodes = make([]*node, 0, len(info.ReadySCs))
	for conn, connInfo := range info.ReadySCs {
		weight := connInfo.Address.Attributes.Value("weight").(int)
		res.nodes = append(res.nodes, &node{
			conn:            conn,
			weight:          weight,
			currentWeight:   weight,
			effectiveWeight: weight,
		})
		res.totalWeight += weight
	}
	return &res
}

type WeightPollingPicker struct {
	mutex       sync.Mutex
	totalWeight int
	nodes       []*node
}

func (w *WeightPollingPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(w.nodes) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	w.mutex.Lock()
	defer w.mutex.Unlock()
	maxWeightNode := w.nodes[0]
	for i := 1; i < len(w.nodes); i++ {
		if cmp.Compare(maxWeightNode.currentWeight, w.nodes[i].currentWeight) < 0 {
			maxWeightNode.currentWeight += maxWeightNode.effectiveWeight
			maxWeightNode = w.nodes[i]
			continue
		}
		w.nodes[i].currentWeight += w.nodes[i].effectiveWeight
	}
	maxWeightNode.currentWeight -= w.totalWeight
	return balancer.PickResult{
		SubConn: maxWeightNode.conn,
		Done: func(info balancer.DoneInfo) {
			//可以考虑有 err 动态调整 effectiveWeight，但可能会导致权重低的节点承担过多流量
		},
	}, nil
}
