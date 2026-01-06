// Package config @Author:冯铁城 [17615007230@163.com] 2025-12-29 19:38:21
package config

import (
	"encoding/json"
	"log"
	"os"
)

// Config 应用配置
type Config struct {
	Server  ServerConfig  `json:"server"`  // 服务配置
	Nacos   NacosConfig   `json:"nacos"`   // Nacos配置
	Message MessageConfig `json:"message"` // 消息配置
	Kafka   KafkaConfig   `json:"kafka"`   // Kafka配置
	Log     LogConfig     `json:"log"`
}

// InitConfig 初始化配置文件
func InitConfig() *Config {

	//1.读取配置文件
	file, err := os.Open("config/properties/config-local.json")
	if err != nil {
		log.Fatalf("[Config] 配置文件读取失败: %v", err)
	}
	defer file.Close()

	//2.解析为结构体
	var config Config
	if err = json.NewDecoder(file).Decode(&config); err != nil {
		log.Fatalf("[Config] 配置解析为JSON失败: %v", err)
	}

	//3.解析服务器配置
	if err = parseServerConfig(&config); err != nil {
		log.Fatalf("[Config] 解析服务器配置失败: %v", err)
	}

	//4.解析Nacos配置
	parseNacosConfig(&config)

	//5.解析消息配置
	parseMessageConfig(&config)

	//6.解析Kafka配置
	parseKafkaConfig(&config)

	//7.解析日志配置
	parseLogConfig(&config)

	//8.返回配置
	log.Printf("[Config] 配置加载成功: [%+v]", config)
	return &config
}
