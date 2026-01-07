// Package config @Author:冯铁城 [17615007230@163.com] 2025-12-29 19:38:21
package config

import (
	_ "embed"
	"encoding/json"

	"github.com/sirupsen/logrus"
)

//go:embed properties/config-local.json
var configData []byte

// Config 应用配置
type Config struct {
	Server   ServerConfig   `json:"server"`    // 服务配置
	Nacos    NacosConfig    `json:"nacos"`     // Nacos配置
	Message  MessageConfig  `json:"message"`   // 消息配置
	Kafka    KafkaConfig    `json:"kafka"`     // Kafka配置
	Log      LogConfig      `json:"log"`       // 日志配置
	WorkPool WorkPoolConfig `json:"work_pool"` // 工作池配置
}

// InitConfig 初始化配置文件
func InitConfig() *Config {

	//1.解析内嵌配置文件
	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		logrus.Fatalf("[Config] 配置解析为JSON失败: %v", err)
	}

	//2.解析服务器配置
	if err := parseServerConfig(&config); err != nil {
		logrus.Errorf("[Config] 解析服务器配置失败: %v", err)
	}

	//3.解析Nacos配置
	parseNacosConfig(&config)

	//4.解析消息配置
	parseMessageConfig(&config)

	//5.解析Kafka配置
	parseKafkaConfig(&config)

	//6.解析日志配置
	parseLogConfig(&config)

	//7.解析工作池配置
	parseWorkPoolConfig(&config)

	//8.返回配置
	logrus.Infof("[Config] 配置加载成功: [%+v]", config)
	return &config
}
