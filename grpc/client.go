package grpc

import (
	"fmt"
	"github.com/msprojectlb/project-common/grpc/registry"
	"github.com/msprojectlb/project-grpc/grpcProject"
	"github.com/msprojectlb/project-grpc/grpcUser"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	UserServiceName    = "project-user"
	ProjectServiceName = "project-project"
)

// NewClient 创建一个grpc客户端
func NewClient(register registry.Registry, balancerName, serviceName string, interceptor ...grpc.UnaryClientInterceptor) (*grpc.ClientConn, error) {
	resolverBuilder := NewResolverBuilder(register)
	target := resolverBuilder.Scheme() + ":///" + serviceName
	return grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(interceptor...),
		grpc.WithResolvers(resolverBuilder),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy":"%s"}`, balancerName)),
	)
}

func ClientWithUser(register registry.Registry, balancerName string, interceptor ...grpc.UnaryClientInterceptor) (grpcUser.UserServiceClient, error) {
	conn, err := NewClient(register, balancerName, UserServiceName, interceptor...)
	if err != nil {
		return nil, err
	}
	return grpcUser.NewUserServiceClient(conn), nil
}

func ClientWithProject(register registry.Registry, balancerName string, interceptor ...grpc.UnaryClientInterceptor) (grpcProject.ProjectServiceClient, error) {
	conn, err := NewClient(register, balancerName, ProjectServiceName, interceptor...)
	if err != nil {
		return nil, err
	}
	return grpcProject.NewProjectServiceClient(conn), nil
}
