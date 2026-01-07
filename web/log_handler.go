// Package web @Author:冯铁城 [17615007230@163.com] 2025-12-29 15:50:02
package web

import (
	"logging-mon-service/commmon/work_pool"
	"logging-mon-service/config"
	"logging-mon-service/model"
	"logging-mon-service/model/res"
	"logging-mon-service/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// logService 日志服务实例
var logService = service.NewLogService()

// uploadLogsAsync 上传日志接口
func uploadLogsAsync(c *gin.Context, cfg *config.Config) {

	//1.声明结构体参数
	var logDto model.LogDto

	//2.获取参数
	if err := c.ShouldBindJSON(&logDto); err != nil {
		logrus.Warnf("[上传日志] 参数校验失败 err:[%v]", err)
		c.JSON(http.StatusBadRequest, res.BadRequestFail(err.Error()))
		return
	}

	//3.请求头获取项目ID以及日志器ID
	projectId := c.GetHeader("X-Project-Id")

	//4.请求头获取日志器ID
	loggerId := c.GetHeader("X-Logger-Id")

	//5.上传日志
	logJob := work_pool.NewLogJob(projectId, loggerId, logDto, cfg)
	work_pool.GlobalLogWorkerPool.Submit(logJob)

	//6.返回
	c.JSON(http.StatusOK, res.CreateSuccess(logDto))
}
