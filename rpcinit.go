package common

import (
	"context"
	"github.com/msprojectlb/project-common/config"
	"github.com/msprojectlb/project-common/grpc"
	"github.com/msprojectlb/project-common/grpc/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type RpcServer struct {
	si           registry.ServiceInstance //grpc服务信息
	registry     registry.Registry        //注册中心
	addr         string                   //grpc服务 监听地址
	registerFunc func(serv *grpc.Server)  //注册rpc服务
	gs           *grpc.Server
}

func NewRpcServer(conf config.GRPCConfig, si registry.ServiceInstance, r registry.Registry) *RpcServer {
	return &RpcServer{
		addr:         conf.Addr,
		registerFunc: conf.RegisterFunc,
		si:           si,
		registry:     r,
		gs:           grpc.NewServer(),
	}
}
func (rs *RpcServer) Start() {
	listen, err := net.Listen("tcp", rs.addr)
	if err != nil {
		log.Fatal(err)
	}
	if rs.registry != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		err := rs.registry.Register(ctx, rs.si)
		cancel()
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			rs.Close()
		}()
	}
	go func() {
		quit := make(chan os.Signal)
		//SIGINT 用户发送INTR字符(Ctrl+C)触发 kill -2
		//SIGTERM 结束程序(可以被捕获、阻塞或忽略)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Printf("grpc server %s Shutting Down at %s... \n", rs.si.Name, rs.si.Address)
		rs.Close()
	}()
	rs.registerFunc(rs.gs)
	rs.gs.Serve(listen)
}
func (rs *RpcServer) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	rs.registry.UnRegister(ctx, rs.si)
	rs.gs.GracefulStop()
}

func ConnectRpc(serviceName string, r registry.Registry) *grpc.ClientConn {
	conn, err := grpc.NewClient("etcd:///"+serviceName, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithResolvers(grpc.NewGrpcResolverBuilder(r)))
	if err != nil {
		log.Fatalf("did not connect %s : %v", serviceName, err)
	}
	return conn
}
