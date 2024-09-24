package byEtcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/msprojectlb/project-common/grpc/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"sync"
)

const PREFIX = "_ETCD_REGISTER_"

type register struct {
	c          *clientv3.Client //etcd 连接
	lease      clientv3.LeaseID //租约
	id         string           //唯一标识
	session    *concurrency.Session
	cancelFunc []func()
	mutex      sync.Mutex
}

// NewRegister c etcd 连接 ttl 租约 秒
func NewRegister(c *clientv3.Client, ttl int) (registry.Registry, error) {
	session, err := concurrency.NewSession(c, concurrency.WithTTL(ttl))
	if err != nil {
		return nil, err
	}
	rs := &register{
		c:          c,
		lease:      session.Lease(),
		id:         uuid.NewString(),
		session:    session,
		cancelFunc: make([]func(), 10),
	}
	return rs, nil
}
func (r *register) Register(ctx context.Context, si registry.ServiceInstance) error {
	marshal, err := json.Marshal(si)
	if err != nil {
		return err
	}
	_, err = r.c.Put(ctx, r.instanceKey(si), string(marshal), clientv3.WithLease(r.lease))
	if err != nil {
		return err
	}
	return nil
}

func (r *register) UnRegister(ctx context.Context, si registry.ServiceInstance) error {
	_, err := r.c.Delete(ctx, r.instanceKey(si))
	return err
}

func (r *register) ListService(ctx context.Context, serviceName string) ([]registry.ServiceInstance, error) {
	response, err := r.c.Get(ctx, r.serviceKey(serviceName), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	service := make([]registry.ServiceInstance, 0, response.Count)
	for _, kv := range response.Kvs {
		s := registry.ServiceInstance{}
		err := json.Unmarshal(kv.Value, &s)
		if err != nil {
			return nil, err
		}
		service = append(service, s)
	}
	return service, nil
}

func (r *register) SubScribe(serviceName string) (<-chan registry.Even, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	r.mutex.Lock()
	if r.cancelFunc != nil {
		r.cancelFunc = append(r.cancelFunc, cancelFunc)
	}
	r.mutex.Unlock()
	ctx = clientv3.WithRequireLeader(ctx)
	watchChan := r.c.Watch(ctx, r.serviceKey(serviceName), clientv3.WithPrefix())
	evenChan := make(chan registry.Even, 10)
	go func() {
		for {
			select {
			case res := <-watchChan:
				if res.Err() != nil {
					return
				}
				if res.Canceled {
					return
				}
				for _, even := range res.Events {
					evenChan <- registry.Even{Type: even.Type.String()}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return evenChan, nil
}

func (r *register) Close() error {
	r.mutex.Lock()
	funcList := r.cancelFunc
	r.cancelFunc = nil
	r.mutex.Unlock()
	for _, f := range funcList {
		f()
	}
	return r.session.Close()
}

func (r *register) instanceKey(si registry.ServiceInstance) string {
	return fmt.Sprintf("%s%s@%s_%s", PREFIX, si.Name, si.Tag, r.id)
}

func (r *register) serviceKey(name string) string {
	return fmt.Sprintf("%s%s", PREFIX, name)
}
