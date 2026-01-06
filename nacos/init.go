// Package nacos @Author:冯铁城 [17615007230@163.com] 2025-12-29 16:30:00
package nacos

import (
	"logging-mon-service/config"

	"github.com/sirupsen/logrus"
)

// Nm 全局Nacos管理器
var Nm *NacosManager

// InitNacosManager 初始化Nacos管理器
func InitNacosManager(c *config.Config) {

	//1.定义异常变量
	var err error

	//2.创建Nacos管理器
	Nm, err = NewNacosManager(c)

	//3.异常不为空，打印日志，终止进程
	if err != nil {
		logrus.Fatalf("[Nacos] 初始化管理器失败: [%v]", err)
	}
}

// RegisterService 注册服务
func RegisterService() {
	if err := Nm.RegisterService(); err != nil {
		logrus.Fatalf("[Nacos] 注册服务失败: [%v]", err)
	}
}

// DeregisterService 注销服务
func DeregisterService() {
	if err := Nm.DeregisterService(); err != nil {
		logrus.Fatalf("[Nacos] 注销服务失败: [%v]", err)
	}
}
