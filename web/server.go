// Package web @Author:冯铁城 [17615007230@163.com] 2025-12-29 15:48:52
package web

import (
	"context"
	"errors"
	"fmt"
	"log"
	"logging-mon-service/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// StartServer 启动服务
func StartServer(c *config.Config) {

	//1.初始化路由
	router := initRouter()

	//2.创建HTTP服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.Server.Port),
		Handler: router,
	}

	//3.协程异步启动服务器
	go func() {
		log.Printf("[Server] 服务启动成功 监听端口: [%d]", c.Server.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Fatalf("[Server] 启动服务器失败: [%v]", err)
		}
	}()

	//4.初始化Nacos管理器
	nm, err := NewNacosManager(c)
	if err != nil {
		log.Fatalf("[Nacos] 初始化管理器失败: [%v]", err)
	}

	//5.注册服务到Nacos
	if err = nm.RegisterService(); err != nil {
		log.Fatalf("[Nacos] 注册服务失败: [%v]", err)
	}

	//6.优雅关闭服务器
	waitForShutdown(server, nm)
}

// waitForShutdown 优雅关闭服务器
func waitForShutdown(server *http.Server, nm *NacosManager) {

	//1.创建信号通道
	quit := make(chan os.Signal, 1)

	//2.监听退出信号
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	//3.阻塞等待信号通道写入退出信号
	<-quit
	log.Println("[Server] 接收到关闭信号，开始优雅关闭...")

	//4.设置关闭超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//5.先从Nacos注销服务
	if err := nm.DeregisterService(); err != nil {
		log.Printf("[Nacos] 注销服务失败: [%v]", err)
	}

	//6.关闭HTTP服务器，等待现有连接完成
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("[Server] 强制关闭服务器: [%v]", err)
	} else {
		log.Println("[Server] 服务器关闭成功")
	}
}
