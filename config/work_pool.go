// Package config @Author:冯铁城 [17615007230@163.com] 2026-01-07 15:47:05
package config

// WorkPoolConfig 工作池配置
type WorkPoolConfig struct {
	Workers int `json:"workers"`  // 工作者数量
	MaxJobs int `json:"max_jobs"` // 最大任务数量
}

// parseWorkPoolConfig 解析工作池配置
func parseWorkPoolConfig(config *Config) {
	if config.WorkPool.Workers <= 0 {
		config.WorkPool.Workers = 10000
	}
	if config.WorkPool.MaxJobs <= 0 {
		config.WorkPool.MaxJobs = 50000
	}
}
