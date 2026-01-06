// Package task @Author:冯铁城 [17615007230@163.com] 2026-01-06 16:25:37
package task

import "logging-mon-service/config"

// iTask 定时任务接口
type iTask interface {
	start()
	stop()
}

// taskManager 任务管理器
type taskManager struct {
	tasks []iTask
}

// newTaskManager 创建任务管理器
func newTaskManager(cfg *config.Config) *taskManager {

	//1.创建任务切片
	tasks := make([]iTask, 0)

	//2.加入任务
	tasks = append(tasks, newResendKafkaMessageTask(cfg))

	//3.创建任务管理器，返回
	return &taskManager{
		tasks: tasks,
	}
}

// startTasks 开始定时任务
func (t *taskManager) startTasks() {
	for _, taskItem := range t.tasks {
		taskItem.start()
	}
}

// stopTasks 停止定时任务
func (t *taskManager) stopTasks() {
	for _, taskItem := range t.tasks {
		taskItem.stop()
	}
}
