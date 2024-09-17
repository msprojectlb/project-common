package mygrpc

import (
	"context"
	"fmt"
	"github.com/msprojectlb/project-common/mygrpc/registry"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type GrpcServer struct {
	si registry.ServiceInstance
	*grpc.Server
	registry registry.Registry
}
type Options func(s *GrpcServer)

func WithRegistry(r registry.Registry) Options {
	return func(s *GrpcServer) {
		s.registry = r
	}
}
func loggerfunc(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	//fmt.Printf("请求Server:%+v\n", info.Server)
	//fmt.Printf("请求FullMethod:%s\n", info.FullMethod)
	res, err := handler(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("请求结束，结果参数为%+v\n", res)
	return res, err
}
func NewGrpcServer(si registry.ServiceInstance, opts ...Options) *GrpcServer {
	rs := &GrpcServer{
		si:     si,
		Server: grpc.NewServer(grpc.UnaryInterceptor(loggerfunc)),
	}
	for _, opt := range opts {
		opt(rs)
	}
	return rs
}

func (s *GrpcServer) Start(addr string) error {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	if s.registry != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		s.si.Address = listen.Addr().String()
		err := s.registry.Register(ctx, s.si)
		cancel()
		if err != nil {
			return err
		}
		defer func() {
			s.Close()
		}()
	}
	go func() {
		quit := make(chan os.Signal)
		//SIGINT 用户发送INTR字符(Ctrl+C)触发 kill -2
		//SIGTERM 结束程序(可以被捕获、阻塞或忽略)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Printf("grpc server %s Shutting Down at %s... \n", s.si.Name, s.si.Address)
		s.Close()
	}()
	err = s.Serve(listen)
	return err
}
func (s *GrpcServer) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	s.registry.UnRegister(ctx, s.si)
	s.GracefulStop()
}
