// Package work_pool @Author:冯铁城 [17615007230@163.com] 2026-01-07 15:58:13
package work_pool

import (
	"logging-mon-service/service"

	"github.com/sirupsen/logrus"
)

// logWorker 日志处理工作协程
type logWorker struct {
	id         int                 // 工作协程ID
	jobs       <-chan LogJob       // 任务队列
	quit       <-chan struct{}     // 退出信号队列
	logService *service.LogService // 日志服务实例
}

// newLogWorker 创建日志处理工作协程
func newLogWorker(id int, jobs <-chan LogJob, quit <-chan struct{}) *logWorker {
	return &logWorker{
		id:         id,
		jobs:       jobs,
		quit:       quit,
		logService: service.NewLogService(),
	}
}

// exe 执行任务
func (r *logWorker) exe() {

	//1.打印启动日志
	logrus.Debugf("[工作池-日志上传] 工作协程[%d]启动", r.id)

	//2.轮询获取任务并执行
	for {
		select {

		//3.从任务队列获取任务
		case job, ok := <-r.jobs:

			//4.如果任务异常，打印日志
			if !ok {
				logrus.Warnf("[工作池-日志上传] 工作协程[%d]任务通道关闭，退出", r.id)
				return
			}

			//5.处理任务
			r.logService.UploadLogs(job.LogDto, job.ProjectId, job.LoggerId, job.Cfg)

		//6.监听退出信号队列，如果收到退出信号，则退出
		case <-r.quit:
			logrus.Debugf("[工作池-日志上传] 工作协程[%d]收到退出信号，退出", r.id)
			return
		}
	}
}
