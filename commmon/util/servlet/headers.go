// Package servlet @Author:冯铁城 [17615007230@163.com] 2026-01-07 10:54:49
package servlet

import "github.com/gin-gonic/gin"

const (
	headerSignature = "X-Signature"
	headerTimestamp = "X-Timestamp"
	headerNonce     = "X-Nonce"
	projectId       = "X-Project-Id"
	loggerId        = "X-Logger-Id"
)

// GetSignature 获取签名
func GetSignature(c *gin.Context) string {
	return c.GetHeader(headerSignature)
}

// GetTimestamp 获取时间戳
func GetTimestamp(c *gin.Context) string {
	return c.GetHeader(headerTimestamp)
}

// GetNonce 获取随机数
func GetNonce(c *gin.Context) string {
	return c.GetHeader(headerNonce)
}

// GetProjectId 获取项目ID
func GetProjectId(c *gin.Context) string {
	return c.GetHeader(projectId)
}

// GetLoggerId 获取日志ID
func GetLoggerId(c *gin.Context) string {
	return c.GetHeader(loggerId)
}
