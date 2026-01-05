// Package web @Author:冯铁城 [17615007230@163.com] 2025-12-29 15:50:02
package web

import (
	"logging-mon-service/model"
	"logging-mon-service/model/res"
	"net/http"

	"github.com/gin-gonic/gin"
)

// uploadLog 上传日志接口
func uploadLog(c *gin.Context) {

	//1.声明结构体参数
	var logDto model.LogDto

	//2.获取参数
	if err := c.ShouldBindJSON(&logDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, res.CreateSuccess(logDto))
	}
}
