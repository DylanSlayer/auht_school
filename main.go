package main

import (
	"auht_school/common"
	"auht_school/routers"
	"auht_school/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func handle1() *gin.Engine {
	r := gin.Default()
	routers.SetRouters(r)
	return r
}
func main() {
	println(utils.GetCurrentTerm())
	common.InitConfig()
	//test()
	//初始化配置
	//设置服务
	server1 := http.Server{
		Addr:         fmt.Sprintf(":%d", common.Settings.ServInfo.Port),
		Handler:      handle1(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println("服务器启动端口：=====>>", common.Settings.ServInfo.Port)
	go func() {
		err := server1.ListenAndServe()
		if err != nil {
			log.Fatalf("listen:%s\n", err)
		}
	}()
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server1.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
