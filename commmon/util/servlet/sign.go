// Package servlet @Author:冯铁城 [17615007230@163.com] 2026-01-07 10:54:39
package servlet

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

// NormalizeRequestBody 规范化请求体
func NormalizeRequestBody(requestBody []byte) string {

	//1.解析JSON为通用对象
	var jsonObject interface{}
	if err := json.Unmarshal(requestBody, &jsonObject); err != nil {
		logrus.Errorf("[拦截器] 签名计算 JSON 解析失败: %v", err)
	}

	//2.重新序列化为规范化字符串
	normalizedBytes, err := json.Marshal(jsonObject)
	if err != nil {
		logrus.Errorf("[拦截器] 签名计算 JSON 序列化失败: %v", err)
	}

	//3.转换为字符串返回
	return string(normalizedBytes)
}

// CalculateSign 计算签名
func CalculateSign(normalizeRequestBody string, nonce string, timestamp string, secretKey string) string {

	//1.获取待签名字符串
	stringToSign := fmt.Sprintf("body=%s&nonce=%s&timestamp=%s", normalizeRequestBody, nonce, timestamp)
	logrus.Infof("[创建签名] 待签名字符串:[%s] 秘钥:[%s]", stringToSign, secretKey)

	//2.使用HMac算法计算签名，返回
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
