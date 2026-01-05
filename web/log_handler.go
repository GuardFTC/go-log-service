// Package web @Author:冯铁城 [17615007230@163.com] 2025-12-29 15:50:02
package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// uploadLog 上传日志接口
func uploadLog(c *gin.Context) {
	fmt.Println("AAA")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}
