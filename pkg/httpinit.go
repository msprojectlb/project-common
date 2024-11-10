package pkg

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/msprojectlb/project-common/logs"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Http struct {
	servName string
	addr     string
	handler  *gin.Engine
	log      *logs.ZapLogger
}

func NewHttp(conf *viper.Viper, handler *gin.Engine, log *logs.ZapLogger) *Http {
	return &Http{
		servName: conf.GetString("app.name"),
		addr:     conf.GetString("app.addr"),
		handler:  handler,
		log:      log,
	}
}
func (h *Http) Run() {
	srv := &http.Server{
		Addr:    h.addr,
		Handler: h.handler,
	}

	go func() {
		defer h.log.Sync()
		h.log.Info("http running...", zap.String("name", h.servName), zap.String("addr", h.addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			h.log.Panic("listen error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	switch sig {
	case syscall.SIGINT, syscall.SIGTERM:
		if err := srv.Shutdown(ctx); err != nil {
			h.log.Error("shutdown server failed", zap.Error(err))
		}
	case syscall.SIGQUIT:
		if err := srv.Close(); err != nil {
			h.log.Error("Close server failed", zap.Error(err))
		}
	}
	h.log.Info("http server exited", zap.String("name", h.servName), zap.String("addr", h.addr), zap.String("signal", sig.String()))
	return
}
