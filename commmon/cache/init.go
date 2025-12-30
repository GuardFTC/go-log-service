// Package cache @Author:冯铁城 [17615007230@163.com] 2025-12-30 15:52:25
package cache

import (
	"logging-mon-service/model"
)

// 全局缓存管理器实例
var logServerCacheManager *LogServerCacheManager

// InitLogServerCache 初始化日志服务缓存
func InitLogServerCache() {

	//1.创建缓存管理器
	logServerCacheManager = NewLogServerCacheManager("logging-mon-server")

	//2.启动定时更新任务
	logServerCacheManager.Start()
}

// StopLogServerCache 停止日志服务缓存
func StopLogServerCache() {
	if logServerCacheManager != nil {
		logServerCacheManager.Stop()
	}
}

// GetProject 获取指定项目
func GetProject(projectID int) *model.Project {

	//1.如果缓存管理器未初始化，则返回nil
	if logServerCacheManager == nil {
		return nil
	}

	//2.获取日志服务对象
	obj := logServerCacheManager.GetLogServerObj()
	if obj == nil {
		return nil
	}

	//3.遍历项目列表，查找指定项目ID的项目
	for _, project := range obj.ProjectObjs {
		if project.ProjectID == projectID {
			return &project
		}
	}

	//4.默认返回
	return nil
}
