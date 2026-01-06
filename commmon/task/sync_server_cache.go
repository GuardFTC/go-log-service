// Package task @Author:冯铁城 [17615007230@163.com] 2026-01-06 15:17:06
package task

import (
	"logging-mon-service/commmon/cache"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

// syncServerCacheTask 同步服务缓存任务
type syncServerCacheTask struct {
	cron *cron.Cron
}

// newSyncServerCacheTask 创建同步服务缓存任务
func newSyncServerCacheTask() *syncServerCacheTask {

	//1.创建任务
	task := &syncServerCacheTask{
		cron: cron.New(cron.WithSeconds()), // 启用秒级 cron
	}

	//2.添加定时任务：30s执行一次
	//执行时间示例
	//2026-01-05 15:00:00
	//2026-01-05 15:00:30
	_, err := task.cron.AddFunc("0/30 * * * * ?", task.syncServerCache)
	if err != nil {
		logrus.Errorf("[定时任务] 同步服务缓存 创建失败: %v", err)
	}

	//3.返回任务实例
	return task
}

// Start 启动定时任务
func (t *syncServerCacheTask) start() {
	t.cron.Start()
	logrus.Infof("[定时任务] 同步服务缓存 已启动，每30s执行一次")
}

// Stop 停止定时任务
func (t *syncServerCacheTask) stop() {
	t.cron.Stop()
	logrus.Infof("[定时任务] 同步服务缓存 已停止")
}

// syncServerCache 同步服务缓存
func (t *syncServerCacheTask) syncServerCache() {
	cache.GLobalServerCacheManager.UpdateCache()
}
