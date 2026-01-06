// Package config @Author:冯铁城 [17615007230@163.com] 2025-12-29 19:35:55
package config

import (
	"bytes"
	"fmt"
	"path"
	"strings"
	"time"

	logrus "github.com/sirupsen/logrus"
)

// LogConfig 日志配置
type LogConfig struct {
	Level  string `json:"level"`
	Color  bool   `json:"color"`
	Format string `json:"format"`
}

// parseLogConfig 解析日志配置
func parseLogConfig(c *Config) {

	//1.根据配置选择格式
	//JSON格式 - 适合生产环境和日志收集
	//自定义文本格式 - 类似Java日志格式
	if c.Log.Format == "json" {
		logrus.SetFormatter(NewJSONFormatter())
	} else {
		logrus.SetFormatter(NewCustomFormatter(c))
	}

	//2.启用调用者信息
	logrus.SetReportCaller(true)

	//3.设置日志级别，默认 Info 及以上输出
	level, _ := logrus.ParseLevel(c.Log.Level)
	logrus.SetLevel(level)
}

// NewJSONFormatter 创建JSON格式器
func NewJSONFormatter() *logrus.JSONFormatter {
	return &logrus.JSONFormatter{
		TimestampFormat: time.DateTime,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "caller",
		},
	}
}

// NewCustomFormatter 创建自定义格式器
func NewCustomFormatter(c *Config) *CustomFormatter {
	return &CustomFormatter{
		ForceColors: c.Log.Color,
	}
}

// CustomFormatter 自定义日志格式器
type CustomFormatter struct {
	ForceColors bool
}

// Format 实现logrus.Formatter接口
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	//1.定义byte缓存
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	//2.时间格式：2026-01-06 11:20:00
	timestamp := entry.Time.Format(time.DateTime)

	//3.日志级别：根据级别添加颜色
	level := strings.ToUpper(entry.Level.String())
	var levelColor string
	if f.ForceColors {
		switch entry.Level {
		case logrus.DebugLevel:
			levelColor = "\033[36m" // 青色
		case logrus.InfoLevel:
			levelColor = "\033[32m" // 绿色
		case logrus.WarnLevel:
			levelColor = "\033[33m" // 黄色
		case logrus.ErrorLevel:
			levelColor = "\033[31m" // 红色
		case logrus.FatalLevel, logrus.PanicLevel:
			levelColor = "\033[35m" // 紫色
		default:
			levelColor = "\033[0m" // 默认
		}
	}

	//3.调用者信息：[package.function(file.go:line)]
	caller := ""
	if entry.HasCaller() {

		//4.获取调用者信息
		funcName := entry.Caller.Function
		fileName := path.Base(entry.Caller.File)
		line := entry.Caller.Line

		//5.提取包名和函数名
		parts := strings.Split(funcName, "/")
		if len(parts) > 0 {
			lastPart := parts[len(parts)-1]
			caller = fmt.Sprintf("[%s(%s:%d)]", lastPart, fileName, line)
		} else {
			caller = fmt.Sprintf("[%s(%s:%d)]", funcName, fileName, line)
		}
	}

	//6.组装日志格式：[时间] - [级别] - [调用者] - 消息
	if f.ForceColors {
		fmt.Fprintf(b, "\033[34m[%s]\033[0m-%s[%s]\033[0m-\033[36m%s\033[0m-%s\n",
			timestamp, levelColor, level, caller, entry.Message)
	} else {
		fmt.Fprintf(b, "[%s]-[%s]-%s-%s\n",
			timestamp, level, caller, entry.Message)
	}

	//7.返回
	return b.Bytes(), nil
}
