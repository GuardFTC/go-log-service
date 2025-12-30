// @Author:冯铁城 [17615007230@163.com] 2025-12-29 15:42:52
package main

import (
	"logging-mon-service/config"
	"logging-mon-service/web"
)

func main() {

	//1.加载配置
	c := config.LoadConfig()

	//2.启动服务器
	web.StartServer(c)
}
