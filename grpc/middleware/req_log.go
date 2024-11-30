package middleware

import (
	"context"
	"github.com/msprojectlb/project-common/logs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

// RecordRequestInfo 记录请求日志
func RecordRequestInfo(log *logs.ZapLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()
		resp, err = handler(ctx, req)
		if err != nil {
			log.Error("rpc resp err",
				zap.String("method", info.FullMethod),
				zap.Duration("time", time.Since(start)),
				zap.Error(err),
				zap.Any("req", req),
				zap.Any("resp", resp),
			)
			return
		}
		log.Info("rpc resp success",
			zap.String("method", info.FullMethod),
			zap.Duration("time", time.Since(start)),
			zap.Any("req", req),
			zap.Any("resp", resp),
		)
		return
	}
}

// RequestDownstreamLogs 记录请求下游日志
func RequestDownstreamLogs(log *logs.ZapLogger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			log.Error("rpc req err",
				zap.String("method", method),
				zap.Duration("time", time.Since(start)),
				zap.Error(err),
				zap.Any("req", req),
				zap.Any("resp", reply),
			)
			return err
		}
		log.Info("rpc req success",
			zap.String("method", method),
			zap.Duration("time", time.Since(start)),
			zap.Any("req", req),
			zap.Any("resp", reply),
		)
		return nil
	}
}
