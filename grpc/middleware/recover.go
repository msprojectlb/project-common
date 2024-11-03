package middleware

import (
	"context"
	"github.com/msprojectlb/project-common/logs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"runtime"
)

// Recovery is a server middleware that recovers from any panics.
func Recovery(log *logs.ZapLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if errInfo := recover(); errInfo != nil {
				buf := make([]byte, 64<<10)
				n := runtime.Stack(buf, false)
				buf = buf[:n]
				log.Error("panic", zap.Any("err", errInfo), zap.Any("req", req), zap.ByteString("stack", buf))
			}
		}()
		return handler(ctx, req)
	}
}
