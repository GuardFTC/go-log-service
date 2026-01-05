// Package model @Author:冯铁城 [17615007230@163.com] 2026-01-05 16:58:45
package model

import (
	"logging-mon-service/model/base"
)

// BaseMessage 基础消息
type BaseMessage struct {
	TableName       string `json:"table_name"`
	TableNameSuffix int    `json:"table_name_suffix"`
}

// LogMessage 日志消息
type LogMessage struct {
	BaseMessage `json:",inline"` //继承基础消息
	ProjectID   int              `json:"project_id"`
	LoggerID    string           `json:"logger_id"`
	Labels      string           `json:"labels"`
	LogLevel    string           `json:"log_level"`
	LogDateTime base.FormatTime  `json:"log_datetime"`
	Content     string           `json:"content"`
}

// NewLogMessage 创建日志消息
func NewLogMessage() *LogMessage {
	return &LogMessage{}
}
