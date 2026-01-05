// Package web @Author:冯铁城 [17615007230@163.com] 2025-12-30 10:31:07
package web

import (
	"logging-mon-service/config"

	"github.com/gin-gonic/gin"
)

// initRouter 初始化路由
func initRouter(c *config.Config) *gin.Engine {

	//1.初始化路由
	router := gin.Default()

	//2.绑定上传日志handler
	logsGroup := router.Group("/api/logs")
	logsGroup.POST("", func(ctx *gin.Context) {
		uploadLogsAsync(ctx, c)
	})

	//3.绑定获取服务列表handler
	nacosGroup := router.Group("/api/nacos")
	nacosGroup.GET("services", getServices)
	nacosGroup.GET("services/:serviceName/instances", getServiceInstances)

	//4.返回
	return router
}
