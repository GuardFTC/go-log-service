// Package config @Author:冯铁城 [17615007230@163.com] 2026-01-05 17:30:00
package config

import (
	"os"
)

// MessageConfig 消息配置
type MessageConfig struct {
	HandlerType string `json:"handler_type"` // 消息处理器类型: kafka_connector, routine_load
}

// parseMessageConfig 解析消息配置
func parseMessageConfig(c *Config) {

	//1.如果环境变量中存在MESSAGE_HANDLER_TYPE，则覆盖消息处理器类型
	if envHandlerType := os.Getenv("MESSAGE_HANDLER_TYPE"); envHandlerType != "" {
		c.Message.HandlerType = envHandlerType
	}

	//2.如果配置为空，设置默认值
	if c.Message.HandlerType == "" {
		c.Message.HandlerType = "kafka_connector"
	}
}
