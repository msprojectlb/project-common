package grpc

import (
	"context"
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
	*grpc.Server                                  // grpc服务
	log             *logs.ZapLogger               // 日志
	registerService func(s grpc.ServiceRegistrar) // 注册服务
	si              registry.ServiceInstance      // 服务实例信息
	registry        registry.Registry             // 注册中心

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

func loggerFunc(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		logs.Helper.Error("req:",
			zap.Any("serv:", info.Server),
			zap.String("method:", info.FullMethod),
			zap.Any("res:", resp),
			zap.Error(err),
		)
		return
	}
	logs.Helper.Info("req:",
		zap.Any("serv:", info.Server),
		zap.String("method:", info.FullMethod),
		zap.Any("res:", resp),
		zap.Error(err),
	)
	return
}
func NewServer(si registry.ServiceInstance, registerService RegisterFunc, opts ...Options) *Server {
	rs := &Server{
		si:              si,
		Server:          grpc.NewServer(grpc.UnaryInterceptor(loggerFunc)),
		log:             logs.Helper,
		registerService: registerService,
	}
	for _, opt := range opts {
		opt(rs)
	}
	if rs.registerService == nil {
		rs.log.LogPanic("微服务尚未注册")
	}
	return rs
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
