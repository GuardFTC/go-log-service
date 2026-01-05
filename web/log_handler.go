// Package web @Author:冯铁城 [17615007230@163.com] 2025-12-29 15:50:02
package web

import (
	"fmt"
	"log"
	"logging-mon-service/commmon/cache"
	"logging-mon-service/commmon/enum"
	"logging-mon-service/commmon/util/message"
	"logging-mon-service/model"
	"logging-mon-service/model/res"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// uploadLogsAsync 上传日志接口
func uploadLogsAsync(c *gin.Context) {

	//1.声明结构体参数
	var logDto model.LogDto

	//2.获取参数
	if err := c.ShouldBindJSON(&logDto); err != nil {
		c.JSON(http.StatusBadRequest, res.BadRequestFail(err.Error()))
		return
	}

	//3.请求头获取项目ID以及日志器ID
	projectId := c.GetHeader("X-Project-Id")
	if projectId == "" {
		c.JSON(http.StatusBadRequest, res.BadRequestFail("项目ID不能为空"))
		return
	}

	//4.请求头获取日志器ID
	loggerId := c.GetHeader("X-Logger-Id")
	if loggerId == "" {
		c.JSON(http.StatusBadRequest, res.BadRequestFail("日志器ID不能为空"))
		return
	}

	//5.上传日志
	uploadLogs(logDto, projectId, loggerId)

	//6.返回
	c.JSON(http.StatusOK, res.CreateSuccess(logDto))
}

// uploadLogs 上传日志
func uploadLogs(logDto model.LogDto, projectId string, loggerId string) {

	//1.获取项目
	project := cache.GetProject(cast.ToInt(projectId))
	if project == nil {
		log.Printf("[上传日志] 项目[%v]不存在", projectId)
		return
	}

	//2.过滤合法日志项
	logItems := make([]model.LogItemDto, 0)
	for _, logItemDto := range logDto.LogItems {

		//3.比较项目允许的日志级别，与日志项的日志级别
		isLegalLogLevel := enum.GreaterThanString(logItemDto.LogLevel, project.LogLevel, true)
		if isLegalLogLevel {

			//4.设置日志器ID
			logItemDto.LoggerID = loggerId

			//5.存入切片
			logItems = append(logItems, logItemDto)
		}
	}

	//6.为空直接返回
	if len(logItems) == 0 {
		log.Printf("[上传日志] 项目[%v]无符合日志级别的日志项", projectId)
		return
	}

	//7.获取消息处理器
	messageHandler := message.Factory.GetMessageHandler(message.RoutineLoad)

	//8.获取Kafka消息
	messages := messageHandler.GetMessages(cast.ToInt(projectId), logItems, 1000)

	//9.为空直接返回
	if len(messages) == 0 {
		log.Printf("[上传日志] 项目[%v]解析Kafka消息为空", projectId)
		return
	}

	fmt.Println(messages)
}
