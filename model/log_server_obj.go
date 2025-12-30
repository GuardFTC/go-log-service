// Package model @Author:冯铁城 [17615007230@163.com] 2025-12-30 15:20:52
package model

// LogServerObj 日志服务对象
type LogServerObj struct {
	ProjectObjs []Project `json:"projectObjs"`
}

// NewLogServerObj 创建日志服务对象
func NewLogServerObj(projectSize int) *LogServerObj {
	return &LogServerObj{
		ProjectObjs: make([]Project, projectSize),
	}
}
