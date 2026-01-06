// Package message @Author:冯铁城 [17615007230@163.com] 2026-01-05 17:05:20
package message

import (
	"logging-mon-service/model"

	"github.com/sirupsen/logrus"
)

// 消息类型常量
const (
	KafkaConnector = "kafka_connector" // Doris Kafka连接器消息
	RoutineLoad    = "routine_load"    // Doris RoutineLoad消息
)

// IMessage 消息接口
type IMessage interface {
	getType() string                                                              // 获取消息类型
	GetMessages(projectId int, logItems []model.LogItemDto, maxSize int) []string // 获取消息
}

// Factory 创建消息处理工厂
var Factory *HandlerFactory

// HandlerFactory 消息处理工厂
type HandlerFactory struct {
	handlerMap map[string]IMessage
}

// InitMessageHandlerFactory 初始化消息处理工厂
func InitMessageHandlerFactory() {

	//1.初始化工厂
	Factory = &HandlerFactory{
		handlerMap: make(map[string]IMessage),
	}

	//2.存入处理器
	Factory.handlerMap[KafkaConnector] = newKafkaConnectorMessage()
	Factory.handlerMap[RoutineLoad] = newRoutineLoadMessage()

	//3.打印日志
	logrus.Infof("[工厂容器] 消息处理工厂初始化完成")
}

// GetMessageHandler 获取消息处理器
func (h *HandlerFactory) GetMessageHandler(handlerType string) IMessage {

	//1.判定是否存在，如果存在，则返回
	if messageHandler, isExist := h.handlerMap[handlerType]; isExist {
		return messageHandler
	}

	//2.否则不存在，打印错误信息并返回nil
	logrus.Errorf("[上传日志] 消息处理器[%v]不存在", handlerType)
	return nil
}
