package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	// 1. 定义输入数据 (使用反引号避免转义地狱)
	requestBody := `{"logItems":[{"logLevel":"info","logDateTime":"2025-12-01 01:23:57","content":"test log content","labels":"[\"info\",\"success\"]"},{"logLevel":"warn","logDateTime":"2025-12-01 06:34:02","content":"test log content","labels":"[\"warn\",\"memory\",\"memory\"]"}]}`
	nonce := "123"
	timeStamp := "1724135987000"
	secretKey := "9CAcMZPf7e3RUY99f4OYnPUqGg6r1gQr"

	// 2. 定义期望值
	expected := "1wTFo1KPybTSys1rq9742atF95NHdTFmaGmjiDqbUg8="

	// 3 & 4. 解析 JSON 为通用对象 (interface{})
	var jsonObject interface{}
	// 将字符串解析为 Go 对象
	if err := json.Unmarshal([]byte(requestBody), &jsonObject); err != nil {
		log.Fatalf("JSON 解析失败: %v", err)
	}

	// 5. 重新序列化为规范化字符串 (Go 的 json.Marshal 默认按 Key 字母排序)
	normalizedBytes, err := json.Marshal(jsonObject)
	if err != nil {
		log.Fatalf("JSON 序列化失败: %v", err)
	}
	normalizeRequestBody := string(normalizedBytes)

	// 6. 获取待签名字符串
	stringToSign := fmt.Sprintf("body=%s&nonce=%s&timestamp=%s", normalizeRequestBody, nonce, timeStamp)

	// 7. 使用 HMac 算法计算签名
	sign := calculateHmacSha256(stringToSign, secretKey)

	// 输出结果
	fmt.Println("--------------------------------------------------")
	fmt.Printf("规范化 Body: %s\n", normalizeRequestBody)
	fmt.Printf("待签名字符:  %s\n", stringToSign)
	fmt.Printf("计算结果:    %s\n", sign)
	fmt.Printf("期望结果:    %s\n", expected)
	fmt.Println("--------------------------------------------------")

	if sign == expected {
		fmt.Println("✅ 验证成功: 签名匹配")
	} else {
		fmt.Println("❌ 验证失败: 签名不匹配")
	}
}

// calculateHmacSha256 计算 HMAC-SHA256 并返回 Base64 字符串
func calculateHmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
