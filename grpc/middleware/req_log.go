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
			log.Error("rpc err",
				zap.String("method", info.FullMethod),
				zap.Error(err),
				zap.Any("req", req),
				zap.Any("resp", resp),
			)
			return
		}
		log.Info("rpc success",
			zap.String("method", info.FullMethod),
			zap.Any("req", req),
			zap.Any("resp", resp),
			zap.Duration("time", time.Since(start)),
		)
		return
	}
}
