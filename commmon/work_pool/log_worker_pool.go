// Package work_pool @Author:冯铁城 [17615007230@163.com] 2026-01-07 15:30:00
package work_pool

import (
	"logging-mon-service/service"
	"sync"

	"github.com/sirupsen/logrus"
)

// LogWorkerPool 日志工作池
type LogWorkerPool struct {
	jobs       chan LogJob
	quit       chan struct{}
	workersNum int
	maxJobsNum int
	logService *service.LogService
	wg         sync.WaitGroup
	mu         sync.RWMutex
	started    bool
}

// NewLogWorkerPool 创建新的日志工作池
func NewLogWorkerPool(workersNum, maxJobsNum int) *LogWorkerPool {
	return &LogWorkerPool{
		jobs:       make(chan LogJob, maxJobsNum),
		quit:       make(chan struct{}),
		workersNum: workersNum,
		maxJobsNum: maxJobsNum,
		logService: service.NewLogService(),
	}
}

// Start 启动工作池
func (p *LogWorkerPool) Start() {

	//1.加锁，确保工作池仅启动一次
	p.mu.Lock()
	defer p.mu.Unlock()

	//2.如果已经启动，返回
	if p.started {
		return
	}

	//3.工作池启动状态设置为已启动
	p.started = true

	//4.启动工作协程
	for i := 0; i < p.workersNum; i++ {

		//5.wg++
		p.wg.Add(1)

		//6.异步创建并启动工作线程
		go func() {
			defer p.wg.Done()
			newLogWorker(i, p.jobs, p.quit).exe()
		}()
	}

	//7.打印启动信息
	logrus.Infof("[工作池-日志上传] 启动完成 工作协程个数:[%v]", p.workersNum)
}

// Stop 停止工作池
func (p *LogWorkerPool) Stop() {

	//1.加锁，确保工作池仅停止一次
	p.mu.Lock()
	defer p.mu.Unlock()

	//2.如果工作池未启动，返回
	if !p.started {
		return
	}

	//3.发送退出信号
	close(p.quit)

	//4.等待所有工作协程完成
	p.wg.Wait()

	//5.所有工作协程都退出后，关闭任务通道
	close(p.jobs)

	//6.工作池状态设置为已停止
	p.started = false

	//7.打印日志
	logrus.Infof("[工作池-日志上传] 停止完成")
}

// Submit 提交任务到工作池
func (p *LogWorkerPool) Submit(job LogJob) bool {

	//1.加锁，确保提交时线程安全
	p.mu.RLock()
	defer p.mu.RUnlock()

	//2.如果工作池未启动，返回
	if !p.started {
		logrus.Warnf("[工作池-日志上传] 工作池未启动，拒绝任务:[%s]", job.String())
		return false
	}

	//3.向任务通道发送任务
	//如果任务通道未关闭，则发送任务，并返回true
	//如果队列已满，则直接处理任务，返回true
	select {
	case p.jobs <- job:
		return true
	default:
		logrus.Warnf("[工作池] 队列已满，当前线程直接处理任务:[%s]", job.String())
		p.logService.UploadLogs(job.LogDto, job.ProjectId, job.LoggerId, job.Cfg)
		return true
	}
}
