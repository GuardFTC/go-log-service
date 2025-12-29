// Package web @Author:冯铁城 [17615007230@163.com] 2025-12-29 15:50:02
package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 全局Nacos管理器
var nacosManager *NacosManager

// uploadLog 上传日志接口
func uploadLog(c *gin.Context) {
	fmt.Println("AAA")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}

// getServices 获取服务列表接口
func getServices(c *gin.Context) {
	if nacosManager == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Nacos管理器未初始化",
		})
		return
	}

	services, err := nacosManager.GetAllServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": fmt.Sprintf("获取服务列表失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    services,
	})
}

// getServiceInstances 获取指定服务的实例列表
func getServiceInstances(c *gin.Context) {
	if nacosManager == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Nacos管理器未初始化",
		})
		return
	}

	serviceName := c.Param("serviceName")
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "服务名称不能为空",
		})
		return
	}

	instances, err := nacosManager.GetServiceInstances(serviceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": fmt.Sprintf("获取服务实例失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    instances,
	})
}
