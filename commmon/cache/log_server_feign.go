// Package task @Author:冯铁城 [17615007230@163.com] 2025-12-30 14:27:29
package cache

import (
	"logging-mon-service/feign"
	"logging-mon-service/model"
)

// LogServerService Server服务客户端
type LogServerService struct {
	client feign.ServiceClient
}

// NewLogServerService 创建Server服务客户端
func NewLogServerService() *LogServerService {
	return &LogServerService{
		client: feign.NewFeignClient("logging-mon-server"),
	}
}

// GetLogServerObj 获取Server服务对象
func (l *LogServerService) GetLogServerObj() (*model.LogServerObj, error) {

	//1.定义result
	var result model.LogServerObj

	//2.Get请求
	if err := l.client.Get("/rpc/log-server/obj", &result); err != nil {
		return nil, err
	}

	//3.返回
	return &result, nil
}
