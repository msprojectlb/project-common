package registry

import (
	"context"
	"io"
)

type Registry interface {
	// Register 将grpc服务器注册到注册中心
	Register(ctx context.Context, si ServiceInstance) error
	UnRegister(ctx context.Context, si ServiceInstance) error

	// ListService 根据serviceName获取可用的服务实例
	ListService(ctx context.Context, serviceName string) ([]ServiceInstance, error)
	SubScribe(serviceName string) (<-chan Even, error)

	io.Closer
}

type ServiceInstance struct {
	Name    string //grpc 服务器名称
	Address string //grpc 服务器地址
	Tag     string //grpc 服务器标签
	Weight  int    //grpc 服务器权重
}

type Even struct {
	Type string
}
