// Package web @Author:冯铁城 [17615007230@163.com] 2025-12-29 15:48:52
package web

import (
	"fmt"
	"log"
	"logging-mon-service/config"

	"github.com/gin-gonic/gin"
)

// StartServer 启动服务
func StartServer() {

	//1.加载配置
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("[Config] 加载配置失败: %v", err)
	}
	log.Printf("[Config] 配置加载成功: [%+v]", c)

	//2.初始化Nacos管理器
	nacosManager, err = NewNacosManager(c)
	if err != nil {
		log.Fatalf("[Nacos] 初始化管理器失败: [%v]", err)
	}

	//3.注册服务到Nacos
	if err = nacosManager.RegisterService(); err != nil {
		log.Fatalf("[Nacos] 注册服务失败: [%v]", err)
	}

	//4.启动优雅关闭监听
	nacosManager.StartGracefulShutdown()

	//5.初始化路由
	router := gin.Default()

	//6.绑定handler
	router.POST("/api/logs", uploadLog)
	router.GET("/api/services", getServices)                                // 获取所有服务列表
	router.GET("/api/services/:serviceName/instances", getServiceInstances) // 获取指定服务实例

	//7.开启服务器
	log.Printf("[Server] 服务启动中...监听端口: [%d]", c.Server.Port)
	if err = router.Run(fmt.Sprintf(":%d", c.Server.Port)); err != nil {
		log.Fatalf("[Server] 启动服务器失败: [%v]", err)
	}
}
