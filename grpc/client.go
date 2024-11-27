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
func NewClient(register registry.Registry, balancerName, serviceName string) (*grpc.ClientConn, error) {
	resolverBuilder := NewResolverBuilder(register)
	target := resolverBuilder.Scheme() + ":///" + serviceName
	return grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(resolverBuilder),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy":"%s"}`, balancerName)),
	)
}

func NewRpcUserClient(register registry.Registry, balancerName string) (grpcUser.UserServiceClient, error) {
	conn, err := NewClient(register, balancerName, UserServiceName)
	if err != nil {
		return nil, err
	}
	return grpcUser.NewUserServiceClient(conn), nil
}

func NewRpcProjectClient(register registry.Registry, balancerName string) (grpcProject.ProjectServiceClient, error) {
	conn, err := NewClient(register, balancerName, UserServiceName)
	if err != nil {
		return nil, err
	}
	return grpcProject.NewProjectServiceClient(conn), nil
}
