// Package work_pool @Author:冯铁城 [17615007230@163.com] 2026-01-07 15:47:05
package work_pool

// GlobalLogWorkerPool 全局日志工作池实例
var GlobalLogWorkerPool *LogWorkerPool

// InitLogWorkerPool 初始化全局日志工作池
func InitLogWorkerPool(workers, maxJobs int) {
	GlobalLogWorkerPool = NewLogWorkerPool(workers, maxJobs)
	GlobalLogWorkerPool.Start()
}

// CloseLogWorkerPool 关闭全局日志工作池
func CloseLogWorkerPool() {
	GlobalLogWorkerPool.Stop()
}
