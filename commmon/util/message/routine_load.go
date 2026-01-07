// Package message @Author:冯铁城 [17615007230@163.com] 2026-01-05 19:18:55
package message

import (
	"encoding/json"
	"fmt"
	"logging-mon-service/model"

	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
)

const (
	defaultMessagePrefix     = "project_logs_%v|"        // 默认消息前缀模版
	defaultMessagePrefixSize = len(defaultMessagePrefix) // 默认消息前缀模版长度
)

// routineLoadMessage RoutineLoad消息
type routineLoadMessage struct{}

// newRoutineLoadMessage 创建RoutineLoad消息
func newRoutineLoadMessage() *routineLoadMessage {
	return &routineLoadMessage{}
}

// GetType 获取消息类型
func (r *routineLoadMessage) getType() string {
	return RoutineLoad
}

// GetMessages 获取消息
func (r *routineLoadMessage) GetMessages(projectId int, logItems []model.LogItemDto, maxSize int) []string {

	//1.日志项转换为LogMessage
	logMessages := make([]*model.LogMessage, 0)
	for _, logItem := range logItems {

		//2.消息转换
		logMessage, err := r.toLogMessage(logItem, projectId)
		if err != nil {
			logrus.Errorf("[上传日志] Doris Routine Load 转换日志失败: %v", err)
			continue
		}

		//3.写入切片
		logMessages = append(logMessages, logMessage)
	}

	//4.JSON序列化
	messageByte, err := json.Marshal(logMessages)
	if err != nil {
		logrus.Errorf("[上传日志] Doris Routine Load 序列化对象失败: %v", err)
		return nil
	}

	//5.判定大小，未达到最大值，直接返回
	if len(messageByte)+defaultMessagePrefixSize < maxSize {
		return r.appendPrefixForMessage(projectId, []string{string(messageByte)})
	}

	//6.否则根据限定大小拆分为多个消息，返回
	return r.appendPrefixForMessage(projectId, r.splitMessageByMaxLength(logMessages, maxSize))
}

// toLogMessage 转换为日志消息
func (r *routineLoadMessage) toLogMessage(logItem model.LogItemDto, projectId int) (*model.LogMessage, error) {

	//1.创建logMessage
	logMessage := model.NewLogMessage()

	//2.复制属性
	if err := copier.Copy(&logMessage, &logItem); err != nil {
		return nil, err
	}

	//3.设置labels
	if labelsStr, err := logItem.Labels.ToJSONString(); err != nil {
		return nil, err
	} else {
		logMessage.Labels = labelsStr
	}

	//4.设置项目ID
	logMessage.ProjectID = projectId

	//5.返回
	return logMessage, nil
}

// splitMessageByMaxLength 按最大长度分割消息
func (r *routineLoadMessage) splitMessageByMaxLength(logMessages []*model.LogMessage, maxSize int) []string {

	//1.定义结果集消息切片
	result := make([]string, 0)

	//2.定义当前批次消息切片
	currentMessages := make([]*model.LogMessage, 0)

	//3.遍历日志消息
	for _, logMessage := range logMessages {

		//4.加入当前批次
		currentMessages = append(currentMessages, logMessage)

		//5.解析为JSON字符串
		currentMessagesByte, _ := json.Marshal(currentMessages)

		//6.检查当前批次大小
		if len(currentMessagesByte)+defaultMessagePrefixSize > maxSize {

			//7.如果只有一个消息，且批次大小超出最大限制
			if len(currentMessages) == 1 {

				//8.打印日志
				logrus.Warnf("[上传日志] Doris Routine Load 单条消息超出最大%v字节限制", maxSize)

				//9.清空当前批次
				currentMessages = make([]*model.LogMessage, 0)

				//10.继续循环
				continue
			}

			//7.如果有多个消息，且批次大小超出最大限制
			if len(currentMessages) > 1 {

				//8.移除当前消息
				currentMessages = currentMessages[:len(currentMessages)-1]

				//9.解析为JSON字符串,加入结果集
				currentMessagesByte, _ := json.Marshal(currentMessages)
				result = append(result, string(currentMessagesByte))

				//10.清空当前批次
				currentMessages = make([]*model.LogMessage, 0)

				//11.加入当前消息
				currentMessages = append(currentMessages, logMessage)
			}
		}
	}

	//12.处理最后一个批次
	if len(currentMessages) > 0 {
		currentMessagesByte, _ := json.Marshal(currentMessages)
		result = append(result, string(currentMessagesByte))
	}

	//13.返回消息结果集
	return result
}

// appendPrefixForMessage 添加消息前缀
func (r *routineLoadMessage) appendPrefixForMessage(projectId int, messages []string) []string {

	//1.定义结果集
	result := make([]string, 0)

	//2.遍历消息
	for _, message := range messages {

		//3.添加消息前缀
		appendPrefixMessage := fmt.Sprintf(defaultMessagePrefix, projectId) + message

		//4.存入结果集
		result = append(result, appendPrefixMessage)
	}

	//5.返回
	return result
}
