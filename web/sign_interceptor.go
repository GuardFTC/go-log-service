// Package web @Author:冯铁城 [17615007230@163.com] 2026-01-07 10:43:11
package web

import (
	"bytes"
	"io"
	"logging-mon-service/commmon/cache"
	"logging-mon-service/commmon/util/servlet"
	"logging-mon-service/model/res"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

// signInterceptor 签名拦截器
func signInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {

		//1.获取请求头中的签名、时间戳、随机数、项目ID信息
		clientSign := servlet.GetSignature(c)
		timestamp := servlet.GetTimestamp(c)
		nonce := servlet.GetNonce(c)
		projectID := servlet.GetProjectId(c)
		loggerID := servlet.GetLoggerId(c)

		//2.判空处理
		if clientSign == "" || timestamp == "" || nonce == "" || projectID == "" || loggerID == "" {
			logrus.Warnf("[拦截器校验] 缺少必要请求头信息 clientSign:[%s] timestamp:[%s] nonce:%s projectId:[%s] loggerId:[%s]", clientSign, timestamp, nonce, projectID, loggerID)
			c.AbortWithStatusJSON(http.StatusBadRequest, res.UnauthorizedFail("缺少必要的请求头参数"))
			return
		}

		//3.根据项目ID获取项目秘钥
		project := cache.GetProject(cast.ToInt(projectID))
		if project == nil {
			logrus.Warnf("[拦截器校验] 项目[%s]不存在", projectID)
			c.AbortWithStatusJSON(http.StatusBadRequest, res.UnauthorizedFail("项目不存在"))
			return
		}

		//4.获取请求体
		body, err := getBody(c)
		if err != nil {
			logrus.Warnf("[拦截器校验] 获取请求体失败 err:[%v]", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, res.UnauthorizedFail("获取请求体失败"))
			return
		}

		//5.将 body 重新赋值给 c.Request.Body，使其可被后续读取
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		//6.按照约定的规则(按字母顺序,对Key进行排序),重新序列化为规范化字符串
		normalizeRequestBody := servlet.NormalizeRequestBody(body)

		//7.计算签名
		sign := servlet.CalculateSign(normalizeRequestBody, nonce, timestamp, project.ProjectKey)

		//8.验证签名是否匹配
		if sign != clientSign {
			logrus.Warnf("[拦截器校验] 签名不匹配 客户端:[%s] 服务端:[%s]", clientSign, sign)
			c.AbortWithStatusJSON(http.StatusBadRequest, res.UnauthorizedFail("签名不匹配"))
			return
		}

		//9.放行请求
		c.Next()
	}
}

// getBody 获取请求体
func getBody(c *gin.Context) ([]byte, error) {

	//1.读取原始 Body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Warnf("[拦截器校验] 获取请求体失败 err:[%v]", err)
		return nil, err
	}

	//2.关闭原始Body流
	if err = c.Request.Body.Close(); err != nil {
		logrus.Warnf("[拦截器校验] 关闭原始Body失败 err:[%v]", err)
		return nil, err
	}

	//3.返回
	return body, nil
}
