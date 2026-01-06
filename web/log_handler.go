// Package web @Author:冯铁城 [17615007230@163.com] 2025-12-29 15:50:02
package web

import (
	"logging-mon-service/commmon/cache"
	"logging-mon-service/commmon/enum"
	"logging-mon-service/commmon/util/message"
	"logging-mon-service/config"
	"logging-mon-service/kafka"
	"logging-mon-service/model"
	"logging-mon-service/model/res"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

// maxMessageSize 最大消息大小
const maxMessageSize = 1024 * 1024

// uploadLogsAsync 上传日志接口
func uploadLogsAsync(c *gin.Context, cfg *config.Config) {

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
	uploadLogs(logDto, projectId, loggerId, cfg)

	//6.返回
	c.JSON(http.StatusOK, res.CreateSuccess(logDto))
}

// uploadLogs 上传日志
func uploadLogs(logDto model.LogDto, projectId string, loggerId string, cfg *config.Config) {

	//1.获取项目
	project := cache.GetProject(cast.ToInt(projectId))
	if project == nil {
		logrus.Errorf("[上传日志] 项目[%v]不存在", projectId)
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
		logrus.Warnf("[上传日志] 项目[%v]无符合日志级别的日志项", projectId)
		return
	}

	//7.获取消息处理器
	messageHandler := message.Factory.GetMessageHandler(cfg.Message.HandlerType)

	//8.获取Kafka消息
	messages := messageHandler.GetMessages(cast.ToInt(projectId), logItems, maxMessageSize)

	//9.为空直接返回
	if len(messages) == 0 {
		logrus.Warnf("[上传日志] 项目[%v]解析Kafka消息为空", projectId)
		return
	}

	//10.判定topic是否为空
	if cfg.Kafka.Topic == "" {
		logrus.Warnf("[上传日志] 项目[%v]Kafka Topic为空", projectId)
		return
	}

	//11.循环发送消息
	for _, kafkaMessage := range messages {

		//12.发送消息
		//发送失败，写入文件
		//发送成功，打印日志
		if err := kafka.GlobalProducer.SendMassage(cfg.Kafka.Topic, kafkaMessage); err != nil {
			if file, e := message.WriteMessageFile(kafkaMessage); e != nil {
				logrus.Errorf("[上传日志] 项目[%v]发送Kafka消息:[%v] 失败:[%v] 写入文件失败:[%v]", projectId, kafkaMessage, err, e)
			} else {
				logrus.Errorf("[上传日志] 项目[%v]发送Kafka消息:[%v] 失败:[%v] 写入文件成功:[%v]", projectId, kafkaMessage, err, file)
			}
		} else {
			logrus.Debugf("[上传日志] 项目[%v]发送Kafka消息:[%v]成功", projectId, kafkaMessage)
		}
	}
}
