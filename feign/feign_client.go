// Package feign 服务间调用客户端
package feign

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"logging-mon-service/nacos"
	"net/http"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/model"
)

// Client 类似OpenFeign的HTTP客户端
type Client struct {
	serviceName  string       // 服务名称
	client       *http.Client // HTTP 客户端
	loadBalancer LoadBalancer // 负载均衡器
}

// NewFeignClient 创建新的Feign客户端
func NewFeignClient(serviceName string) *Client {
	return &Client{
		serviceName:  serviceName,
		loadBalancer: GetLoadBalancer(RoundRobin),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewFeignClientWithLoadBalancer 创建带指定负载均衡策略的Feign客户端
func NewFeignClientWithLoadBalancer(serviceName string, lbType LoadBalancerType) *Client {
	return &Client{
		serviceName:  serviceName,
		loadBalancer: GetLoadBalancer(lbType),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// -------------------------------------------------------外部方法-------------------------------------------//

// Get 发送GET请求
func (f *Client) Get(path string, result interface{}) error {
	return f.request("GET", path, nil, result)
}

// Post 发送POST请求
func (f *Client) Post(path string, body interface{}, result interface{}) error {
	return f.request("POST", path, body, result)
}

// Put 发送PUT请求
func (f *Client) Put(path string, body interface{}, result interface{}) error {
	return f.request("PUT", path, body, result)
}

// Delete 发送DELETE请求
func (f *Client) Delete(path string, result interface{}) error {
	return f.request("DELETE", path, nil, result)
}

// -------------------------------------------------------内部方法-------------------------------------------//

// request 执行HTTP请求，包含负载均衡逻辑
func (f *Client) request(method, path string, body interface{}, result interface{}) error {

	//1.获取可用且健康服务实例列表
	instances, err := f.getHealthInstances()
	if err != nil {
		return err
	}

	//2.负载均衡选择实例
	instance := f.loadBalancer.Select(instances)

	//3.构建请求
	req, err := f.buildRequest(method, path, body, instance)
	if err != nil {
		return err
	}

	//4.发送请求
	resp, err := f.client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败:[%v]", err)
	}
	defer resp.Body.Close()

	//5.处理响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败:[%v]", err)
	}

	//6.如果响应状态码不为2XX，返回错误
	if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusIMUsed) {
		return fmt.Errorf("请求失败,状态码:[%d],响应:[%s]", resp.StatusCode, string(respBody))
	}

	//7.如果result不为空，尝试将响应数据解析为result
	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("解析响应失败:[%v]", err)
		}
	}

	//8.默认返回
	return nil
}

// getHealthInstances 获取可用且健康服务实例列表
func (f *Client) getHealthInstances() ([]model.Instance, error) {

	//1.从Nacos获取服务实例
	instances, err := nacos.Nm.GetServiceInstances(f.serviceName)
	if err != nil {
		return nil, fmt.Errorf("获取服务实例失败: %v", err)
	}

	//2.如果没有可用实例，返回错误
	if len(instances) == 0 {
		return nil, fmt.Errorf("[%s]获取服务实例列表为空", f.serviceName)
	}

	//3.过滤出可用且健康的服务实例
	var healthyInstances []model.Instance
	for _, instance := range instances {

		//4.如果服务可用且健康，加入结果列表
		if instance.Enable && instance.Healthy {
			healthyInstances = append(healthyInstances, instance)
		}
	}

	//5.如果没有可用实例，返回错误
	if len(healthyInstances) == 0 {
		return nil, fmt.Errorf("[%s]无可用实例", f.serviceName)
	}

	//6.返回
	return instances, nil
}

// buildRequest 构建HTTP请求
func (f *Client) buildRequest(method string, path string, body interface{}, instance model.Instance) (*http.Request, error) {

	//1.拼接请求URL
	url := fmt.Sprintf("http://%s:%d%s", instance.Ip, instance.Port, path)

	//2.序列化请求体
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	//3.创建请求
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	//4.设置请求类型
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	//5.返回请求
	return req, nil
}
