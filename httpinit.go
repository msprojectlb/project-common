package common

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpInit struct {
	servName string
	addr     string
	handler  *gin.Engine
}

func NewHttpInit(conf *viper.Viper, handler *gin.Engine) *HttpInit {
	return &HttpInit{
		servName: conf.GetString("server.name"),
		addr:     conf.GetString("server.addr"),
		handler:  handler,
	}
}
func (h *HttpInit) Run() {
	srv := &http.Server{
		Addr:    h.addr,
		Handler: h.handler,
	}
	//保证下面的优雅启停
	go func() {
		log.Printf("%s running in %s \n", h.servName, srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal)
	//SIGINT 用户发送INTR字符(Ctrl+C)触发 kill -2
	//SIGTERM 结束程序(可以被捕获、阻塞或忽略)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutting Down project %s... \n", h.servName)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("%s Shutdown, cause by : %v", h.servName, err)
	}
	select {
	case <-ctx.Done():
		log.Println("wait timeout....")
	}
	log.Printf("%s stop success... \n", h.servName)
}
