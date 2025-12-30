// Package config @Author:冯铁城 [17615007230@163.com] 2025-12-29 19:36:57
package config

import (
	"os"
)

// NacosConfig Nacos配置
type NacosConfig struct {
	ServerAddr string `json:"server_addr"` // Nacos地址
	ServerPort uint64 `json:"server_port"` // Nacos端口
	Namespace  string `json:"namespace"`   // 命名空间
	Group      string `json:"group"`       // 组名
	Timeout    uint64 `json:"timeout"`     // 超时时间
	LogDir     string `json:"log_dir"`     // 日志目录
	CacheDir   string `json:"cache_dir"`   // 缓存目录
}

// parseNacosConfig 加载服务配置
func parseNacosConfig(c *Config) {

	//1.如果环境变量中存在NACOS_ADDR，则覆盖Nacos地址
	if envNacosAddr := os.Getenv("NACOS_ADDR"); envNacosAddr != "" {
		c.Nacos.ServerAddr = envNacosAddr
	}
}
