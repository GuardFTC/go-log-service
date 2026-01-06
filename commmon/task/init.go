// Package task @Author:冯铁城 [17615007230@163.com] 2026-01-06 16:39:14
package task

import "logging-mon-service/config"

// 全局任务管理器
var globalTaskManager *taskManager

// InitTaskManager 初始化任务管理器
func InitTaskManager(cfg *config.Config) {

	//1.创建任务管理器
	globalTaskManager = newTaskManager(cfg)

	//2.启动任务
	globalTaskManager.startTasks()
}

// StopTaskManager 停止任务管理器
func StopTaskManager() {
	if globalTaskManager != nil {
		globalTaskManager.stopTasks()
	}
}
