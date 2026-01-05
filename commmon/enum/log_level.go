// Package enum @Author:冯铁城 [17615007230@163.com] 2026-01-05 16:31:05
package enum

import (
	"strings"
)

// LogLevel 定义日志级别类型
type LogLevel int

// 常量定义
const (
	Trace LogLevel = iota + 1
	Debug
	Info
	Warn
	Error
)

// String 返回日志级别的字符串表示
func (l LogLevel) String() string {
	switch l {
	case Trace:
		return "trace"
	case Debug:
		return "debug"
	case Info:
		return "info"
	case Warn:
		return "warn"
	case Error:
		return "error"
	default:
		return "trace"
	}
}

// ParseLogLevel 从字符串解析为 LogLevel，支持别名（如 warning -> warn, err -> error）
func ParseLogLevel(s string) LogLevel {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return Debug
	case "info":
		return Info
	case "warn", "warning":
		return Warn
	case "error", "err":
		return Error
	default:
		return Trace
	}
}

// AllValues 返回所有日志级别字符串列表
func AllValues() []string {
	return []string{"trace", "debug", "info", "warn", "error"}
}

// GreaterThanString 判断日志级别是否大于等于某个日志级别
func GreaterThanString(left string, right string, isEqual bool) bool {
	return GreaterThanEnumAndString(ParseLogLevel(left), right, isEqual)
}

// GreaterThanEnumAndString 判断日志级别是否大于等于某个日志级别
func GreaterThanEnumAndString(left LogLevel, right string, isEqual bool) bool {
	return GreaterThanEnum(left, ParseLogLevel(right), isEqual)
}

// GreaterThanEnum 判断日志级别是否大于等于某个日志级别
func GreaterThanEnum(left, right LogLevel, isEqual bool) bool {
	if isEqual {
		return left >= right
	}
	return left > right
}
