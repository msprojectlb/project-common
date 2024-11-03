package grpc

import (
	"context"
	"github.com/msprojectlb/project-common/grpc/middleware"
	"github.com/msprojectlb/project-common/grpc/registry"
	"github.com/msprojectlb/project-common/logs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	*grpc.Server                                   // grpc服务
	log              *logs.ZapLogger               // 日志
	registerService  func(s grpc.ServiceRegistrar) // 注册服务
	si               registry.ServiceInstance      // 服务实例信息
	registry         registry.Registry             // 注册中心
	unaryInterceptor []grpc.UnaryServerInterceptor // 中间件
}

type Options func(s *Server)

type RegisterFunc func(s grpc.ServiceRegistrar)

// WithRegistry 注入注册中心
func WithRegistry(r registry.Registry) Options {
	return func(s *Server) {
		s.registry = r
	}
}

// WithLogger 注入日志
func WithLogger(log *logs.ZapLogger) Options {
	return func(s *Server) {
		s.log = log
	}
}

// WithUnaryInterceptor 注入中间件
func WithUnaryInterceptor(middleware ...grpc.UnaryServerInterceptor) Options {
	return func(s *Server) {
		s.unaryInterceptor = append(s.unaryInterceptor, middleware...)
	}
}

func NewServer(si registry.ServiceInstance, registerService RegisterFunc, opts ...Options) *Server {
	s := &Server{
		si:              si,
		log:             logs.Helper,
		registerService: registerService,
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.registerService == nil {
		s.log.Panic("微服务尚未注册")
	}

	middleWares := make([]grpc.UnaryServerInterceptor, 0, len(s.unaryInterceptor)+2)
	// recover中间件
	middleWares = append(middleWares, middleware.Recovery(s.log))
	// 请求日志中间件
	middleWares = append(middleWares, middleware.RecordRequestInfo(s.log))
	middleWares = append(middleWares, s.unaryInterceptor...)
	s.Server = grpc.NewServer(grpc.ChainUnaryInterceptor(middleWares...))

	return s
}

func (s *Server) Start() {
	listen, err := net.Listen("tcp", s.si.Address)
	if err != nil {
		s.log.Panic("net.listen error", zap.Error(err))
	}
	if s.registry != nil {
		s.si.Address = listen.Addr().String()
		err = s.registry.Register(context.Background(), s.si)
		if err != nil {
			err = listen.Close()
			if err != nil {
				s.log.Panic("listen.Close error", zap.Error(err))
			}
			s.log.Panic("registry.Register error", zap.Error(err))
		}
	}
	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		sig := <-quit
		s.Close(sig)
	}()
	s.registerService(s)
	if err = s.Serve(listen); err != nil {
		s.log.Panic("grpc.Serve error", zap.Error(err))
	}
}
func (s *Server) Close(sig os.Signal) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := s.registry.UnRegister(ctx, s.si)
	if err != nil {
		s.log.Error("registry.UnRegister error", zap.Error(err))
	}
	s.log.Info("server shutdown", zap.String("signal", sig.String()))
	s.log.Sync()
	switch sig {
	case syscall.SIGINT, syscall.SIGTERM:
		s.GracefulStop()
	case syscall.SIGQUIT:
		s.Stop()
	}
}
