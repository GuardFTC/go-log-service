// @Author:冯铁城 [17615007230@163.com] 2025-12-29 15:42:52
package main

import (
	"logging-mon-service/commmon/cache"
	"logging-mon-service/commmon/util/message"
	"logging-mon-service/config"
	"logging-mon-service/kafka"
	"logging-mon-service/nacos"
	"logging-mon-service/web"
)

func main() {

	//1.加载配置
	c := config.InitConfig()

	//2.初始化Nacos管理器
	nacos.InitNacosManager(c)

	//3.初始化日志服务缓存
	cache.InitLogServerCache()
	defer cache.StopLogServerCache()

	//4.初始化消息处理器工厂
	message.InitMessageHandlerFactory()

	//5.初始化Kafka生产者
	kafka.InitProducer(c)
	defer kafka.CloseProducer()

	//6.启动服务器
	web.StartServer(c)
}
