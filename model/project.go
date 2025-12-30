// Package model @Author:冯铁城 [17615007230@163.com] 2025-12-30 14:22:00
package model

// Project 项目对象
type Project struct {
	ProjectID   int    `json:"projectId,omitempty"`   // 项目ID
	ProjectName string `json:"projectName,omitempty"` // 项目名称
	ProjectKey  string `json:"projectKey,omitempty"`  // 项目Key
	Remark      string `json:"remark,omitempty"`      // 项目描述
	LogLevel    string `json:"logLevel,omitempty"`    // 项目日志
}
