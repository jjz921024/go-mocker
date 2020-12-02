package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go_mocker/controller"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	r := gin.Default()
	r.Any("/mirror", controller.Mirror)
	r.Any("/monitor", controller.WsUpgradeHandler)

	// 不使用Run方法的目的是要获取到http.Server实例
	//r.Run(":8080")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// 在协程里监听端口，不要阻塞
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %s\n", err)
		}
	}()

	/*go func() {
	  for true {
	    time.Sleep(3 * time.Second)
	    log.Printf("---: %d\n", len(service.ConnSet))
	  }
	}()*/

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 用系统信号阻塞
	<-quit
	log.Println("start shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("sever shutdown error: ", err)
	}
	log.Println("serve exit")

}
