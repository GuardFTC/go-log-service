// Package web @Author:冯铁城 [17615007230@163.com] 2025-12-29 15:50:02
package web

import (
	"fmt"
	"logging-mon-service/model/res"
	"logging-mon-service/nacos"
	"net/http"

	"github.com/gin-gonic/gin"
)

// getServices 获取服务列表接口
func getServices(c *gin.Context) {

	//1.如果Nacos管理器未初始化，返回错误
	if nacos.Nm == nil {
		c.JSON(http.StatusInternalServerError, res.ServerFail("Nacos管理器未初始化").ToJson())
		return
	}

	//2.获取全部服务
	services, err := nacos.Nm.GetAllServices()

	//3.如果异常不为空，返回错误
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.ServerFail(fmt.Sprintf("获取服务列表失败: %v", err)).ToJson())
		return
	}

	//4.否则返回成功
	c.JSON(http.StatusOK, res.QuerySuccess(services))
}

// getServiceInstances 获取指定服务的实例列表
func getServiceInstances(c *gin.Context) {

	//1.如果Nacos管理器未初始化，返回错误
	if nacos.Nm == nil {
		c.JSON(http.StatusInternalServerError, res.ServerFail("Nacos管理器未初始化").ToJson())
		return
	}

	//2.获取服务名称
	serviceName := c.Param("serviceName")

	//3.服务名称为空，返回错误
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, res.BadRequestFail("服务名称不能为空").ToJson())
		return
	}

	//4.根据服务名称获取服务实例
	instances, err := nacos.Nm.GetServiceInstances(serviceName)

	//5.如果异常不为空，返回错误
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.ServerFail(fmt.Sprintf("获取服务实例列表失败: %v", err)))
		return
	}

	//6.否则返回成功
	c.JSON(http.StatusOK, res.QuerySuccess(instances))
}
