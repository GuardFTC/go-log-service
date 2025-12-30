// Package web @Author:冯铁城 [17615007230@163.com] 2025-12-30 10:31:07
package web

import "github.com/gin-gonic/gin"

// initRouter 初始化路由
func initRouter() *gin.Engine {

	//1.初始化路由
	router := gin.Default()

	//2.绑定handler
	router.POST("/api/logs", uploadLog)
	router.GET("/api/services", getServices)                                // 获取所有服务列表
	router.GET("/api/services/:serviceName/instances", getServiceInstances) // 获取指定服务实例

	//3.返回
	return router
}
