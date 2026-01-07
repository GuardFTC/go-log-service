// @Author:冯铁城 [17615007230@163.com] 2025-12-29 15:42:52
package main

import (
	"logging-mon-service/commmon/cache"
	"logging-mon-service/commmon/task"
	"logging-mon-service/commmon/util/message"
	"logging-mon-service/commmon/work_pool"
	"logging-mon-service/config"
	"logging-mon-service/kafka"
	"logging-mon-service/nacos"
	"logging-mon-service/web"
)

func main() {

	//1.加载配置
	c := config.InitConfig()

	//2.初始化消息处理器工厂
	message.InitMessageHandlerFactory()

	//3.初始化日志服务缓存
	cache.InitLogServerCache()

	//4.初始化Kafka生产者
	kafka.InitProducer(c)
	defer kafka.CloseProducer()

	//5.初始化Nacos管理器
	nacos.InitNacosManager(c)
	defer nacos.DeregisterService()

	//6.初始化定时任务
	task.InitTaskManager(c)
	defer task.StopTaskManager()

	//7.初始化日志工作池
	work_pool.InitLogWorkerPool(c.WorkPool.Workers, c.WorkPool.MaxJobs)
	defer work_pool.CloseLogWorkerPool()

	//8.启动服务器
	web.StartServer(c)
}
