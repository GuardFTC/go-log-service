// Package model @Author:冯铁城 [17615007230@163.com] 2026-01-05 15:13:14
package model

import (
	"encoding/json"
	"logging-mon-service/model/base"
)

// LogItemDto 日志项
type LogItemDto struct {
	LoggerID    string               `json:"loggerId" binding:"omitempty"`
	Labels      base.StringJSONArray `json:"labels" binding:"required"`
	LogLevel    string               `json:"logLevel" binding:"required"`
	Content     string               `json:"content" binding:"required"`
	LogDateTime *base.FormatTime     `json:"logDateTime" binding:"required"`
}

// LogDto 日志
type LogDto struct {
	LogItems []LogItemDto `json:"logItems" binding:"required,dive"`
}

// String 日志转字符串
func (l *LogDto) String() string {
	b, _ := json.Marshal(l)
	return string(b)
}
