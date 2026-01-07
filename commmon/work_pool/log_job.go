// Package work_pool @Author:冯铁城 [17615007230@163.com] 2026-01-07 16:48:06
package work_pool

import (
	"logging-mon-service/config"
	"logging-mon-service/model"
)

// LogJob 日志处理任务
type LogJob struct {
	ProjectId string
	LoggerId  string
	LogDto    model.LogDto
	Cfg       *config.Config
}

// NewLogJob 创建日志处理任务
func NewLogJob(projectId string, loggerId string, logDto model.LogDto, cfg *config.Config) *LogJob {
	return &LogJob{
		ProjectId: projectId,
		LoggerId:  loggerId,
		LogDto:    logDto,
		Cfg:       cfg,
	}
}

// String 日志任务转字符串
func (p *LogJob) String() string {
	return "项目ID:" + p.ProjectId + " 日志ID:" + p.LoggerId + " 日志内容:" + p.LogDto.String()
}
