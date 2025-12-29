// Package config @Author:冯铁城 [17615007230@163.com] 2025-12-29 19:35:55
package config

import (
	"net"
	"os"
	"strconv"
)

// ServerConfig 服务配置
type ServerConfig struct {
	Name    string `json:"name"`    // 服务名称
	IP      string `json:"ip"`      // 服务IP
	Port    uint64 `json:"port"`    // 服务端口
	Version string `json:"version"` // 服务版本
}

// parseServerConfig 加载服务配置
func parseServerConfig(c *Config) error {

	//1.如果IP为空
	if c.Server.IP == "" {

		//2.获取本机IP
		ip, err := getLocalIP()
		if err != nil {
			return err
		}

		//3.设置服务IP
		c.Server.IP = ip
	}

	//4.如果环境变量中存在SERVICE_IP和SERVICE_PORT，则覆盖服务IP和端口
	if envIP := os.Getenv("SERVICE_IP"); envIP != "" {
		c.Server.IP = envIP
	}
	if envPort := os.Getenv("SERVICE_PORT"); envPort != "" {
		if port, err := strconv.ParseUint(envPort, 10, 64); err == nil {
			c.Server.Port = port
		}
	}

	//5.返回
	return nil
}

// getLocalIP 获取本机IP地址
func getLocalIP() (string, error) {

	//1.链接本机
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	//2.获取本机IP，返回
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}
