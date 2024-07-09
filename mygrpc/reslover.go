package mygrpc

import (
	"context"
	"github.com/msprojectlb/project-common/mygrpc/registry"
	"google.golang.org/grpc/resolver"
	"time"
)

const SCHEME = "etcd"

type grpcResolverBuilder struct {
	r registry.Registry
}

func NewGrpcResolverBuilder(r registry.Registry) *grpcResolverBuilder {
	return &grpcResolverBuilder{
		r: r,
	}
}

func (b *grpcResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &Resolver{
		cc:     cc,
		r:      b.r,
		target: target,
		close:  make(chan struct{}),
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	go r.watcher()
	return r, nil
}

func (b *grpcResolverBuilder) Scheme() string {
	return SCHEME
}

type Resolver struct {
	cc     resolver.ClientConn
	r      registry.Registry
	target resolver.Target
	close  chan struct{}
}

func (r *Resolver) ResolveNow(options resolver.ResolveNowOptions) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelFunc()
	serviceInstances, err := r.r.ListService(ctx, r.target.Endpoint())
	if err != nil {
		r.cc.ReportError(err)
		return
	}
	address := make([]resolver.Address, 0, len(serviceInstances))
	for _, instance := range serviceInstances {
		address = append(address, resolver.Address{Addr: instance.Address, ServerName: instance.Name})
	}
	err = r.cc.UpdateState(resolver.State{
		Addresses: address,
	})
	//添加延迟时间，否则自定义轮询插件会拿不到可用连接
	time.Sleep(time.Millisecond * 10)
	if err != nil {
		r.cc.ReportError(err)
		return
	}
}
func (r *Resolver) watcher() {
	evens, err := r.r.SubScribe(r.target.Endpoint())
	if err != nil {
		r.cc.ReportError(err)
		return
	}
	for {
		select {
		case <-evens:
			//TODO 根据even 采取del update等操作
			r.ResolveNow(resolver.ResolveNowOptions{})
		case <-r.close:
			return
		}
	}
}
func (r *Resolver) Close() {
	close(r.close)
}
