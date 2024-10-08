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
	si registry.ServiceInstance
	*grpc.Server
	registry registry.Registry
	log      *logs.ZapLogger
}
type Options func(s *Server)

// WithRegistry 注入注册中心
func WithRegistry(r registry.Registry) Options {
	return func(s *Server) {
		s.registry = r
	}
}

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
	}
	logs.Helper.Info("req:",
		zap.Any("serv:", info.Server),
		zap.String("method:", info.FullMethod),
		zap.Any("res:", resp),
		zap.Error(err),
	)
	return
}
func NewServer(si registry.ServiceInstance, opts ...Options) *Server {
	rs := &Server{
		si:     si,
		Server: grpc.NewServer(grpc.UnaryInterceptor(loggerFunc)),
	}
	for _, opt := range opts {
		opt(rs)
	}
	return rs
}

func (s *Server) Start() {
	listen, err := net.Listen("tcp", s.si.Address)
	if err != nil {
		s.log.Panic("net.listen error", zap.Error(err))
	}
	if s.registry != nil {
		defer s.Close()
		s.si.Address = listen.Addr().String()
		err = s.registry.Register(context.Background(), s.si)
		if err != nil {
			s.log.Panic("registry.Register error", zap.Error(err))
		}
	}
	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		s.log.Info("grpc server Shutting Down", zap.String("name", s.si.Name), zap.String("addr", s.si.Address))
		s.Close()
	}()
	if err = s.Serve(listen); err != nil {
		s.log.Panic("grpc.Serve error", zap.Error(err))
	}
}
func (s *Server) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defer s.log.Sync()
	err := s.registry.UnRegister(ctx, s.si)
	if err != nil {
		s.log.Error("registry.UnRegister error", zap.Error(err))
	}
	s.GracefulStop()
}
