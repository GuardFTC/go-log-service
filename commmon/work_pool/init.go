// Package work_pool @Author:冯铁城 [17615007230@163.com] 2026-01-07 15:47:05
package work_pool

import "github.com/sirupsen/logrus"

// GlobalLogWorkerPool 全局日志工作池实例
var GlobalLogWorkerPool *LogWorkerPool

// InitLogWorkerPool 初始化全局日志工作池
func InitLogWorkerPool(workers, maxJobs int) {
	GlobalLogWorkerPool = NewLogWorkerPool(workers, maxJobs)
	GlobalLogWorkerPool.Start()
	logrus.Infof("[工作池] 日志工作池初始化完成 工作协程数:[%d] 最大队列长度:[%d]", workers, maxJobs)
}
