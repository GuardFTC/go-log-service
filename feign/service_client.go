// Package feign @Author:冯铁城 [17615007230@163.com] 2025-12-30 11:23:15
package feign

// ServiceClient 服务客户端接口
type ServiceClient interface {
	Get(path string, result interface{}) error
	Post(path string, body interface{}, result interface{}) error
	Put(path string, body interface{}, result interface{}) error
	Delete(path string, result interface{}) error
}
