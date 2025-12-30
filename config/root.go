// Package config @Author:冯铁城 [17615007230@163.com] 2025-12-29 19:38:21
package config

import (
	"encoding/json"
	"log"
	"os"
)

// Config 应用配置
type Config struct {
	Server ServerConfig `json:"server"` // 服务配置
	Nacos  NacosConfig  `json:"nacos"`  // Nacos配置
}

// LoadConfig 加载配置文件
func LoadConfig() *Config {

	//1.读取配置文件
	file, err := os.Open("config/properties/config-local.json")
	if err != nil {
		log.Fatalf("[Config] 加载配置失败: %v", err)
	}
	defer file.Close()

	//2.解析为结构体
	var config Config
	if err = json.NewDecoder(file).Decode(&config); err != nil {
		log.Fatalf("[Config] 加载配置失败: %v", err)
	}

	//3.解析服务器配置
	if err = parseServerConfig(&config); err != nil {
		log.Fatalf("[Config] 加载配置失败: %v", err)
	}

	//4.解析Nacos配置
	parseNacosConfig(&config)

	//5.返回配置
	log.Printf("[Config] 配置加载成功: [%+v]", config)
	return &config
}
