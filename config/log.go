// Package config @Author:冯铁城 [17615007230@163.com] 2025-12-29 19:35:55
package config

import logrus "github.com/sirupsen/logrus"

// LogConfig 日志配置
type LogConfig struct {
	Level string `json:"level"`
	Color bool   `json:"color"`
}

// parseLogConfig 解析日志配置
func parseLogConfig(c *Config) {

	//1.设置日志格式带颜色
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   c.Log.Color,
		FullTimestamp: true,
	})

	//2.设置日志级别，默认 Info 及以上输出
	level, _ := logrus.ParseLevel(c.Log.Level)
	logrus.SetLevel(level)
}
