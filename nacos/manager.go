// Package nacos @Author:冯铁城 [17615007230@163.com] 2025-12-30 10:58:45
package nacos

import (
	"fmt"
	"log"
	"logging-mon-service/config"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// NacosManager Nacos管理器
type NacosManager struct {
	client naming_client.INamingClient
	config *config.Config
}

// NewNacosManager 创建Nacos管理器
func NewNacosManager(config *config.Config) (*NacosManager, error) {

	//1.创建服务配置
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: config.Nacos.ServerAddr,
			Port:   config.Nacos.ServerPort,
		},
	}

	//2.创建Nacos配置
	clientConfig := constant.ClientConfig{
		NamespaceId:         config.Nacos.Namespace,
		TimeoutMs:           config.Nacos.Timeout,
		NotLoadCacheAtStart: true,
		LogDir:              config.Nacos.LogDir,
		CacheDir:            config.Nacos.CacheDir,
		LogLevel:            "info",
	}

	//3.创建客户端
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return nil, err
	}

	//4.创建Nacos管理器,返回
	return &NacosManager{
		client: client,
		config: config,
	}, nil
}

// RegisterService 注册服务
func (nm *NacosManager) RegisterService() error {

	//1.注册服务
	success, err := nm.client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          nm.config.Server.IP,
		Port:        nm.config.Server.Port,
		ServiceName: nm.config.Server.Name,
		GroupName:   nm.config.Nacos.Group,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata: map[string]string{
			"version": nm.config.Server.Version,
		},
	})

	//2.异常判定
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("注册API返回[false]")
	}

	//3.打印日志
	log.Printf("[Nacos] 服务注册成功: ip:[%s] port:[%d] serviceName:[%s] version:[%s]", nm.config.Server.IP, nm.config.Server.Port, nm.config.Server.Name, nm.config.Server.Version)

	//4.默认返回空异常
	return nil
}

// DeregisterService 注销服务
func (nm *NacosManager) DeregisterService() error {

	//1.注销服务
	success, err := nm.client.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          nm.config.Server.IP,
		Port:        nm.config.Server.Port,
		ServiceName: nm.config.Server.Name,
		GroupName:   nm.config.Nacos.Group,
		Ephemeral:   true,
	})

	//2.异常判定
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("注销API返回[false]")
	}

	//3.打印日志
	log.Printf("[Nacos] 服务注销成功: ip=[%s] port=[%d] serviceName=[%s] version=[%s]", nm.config.Server.IP, nm.config.Server.Port, nm.config.Server.Name, nm.config.Server.Version)

	//4.默认返回空异常
	return nil
}

// GetAllServices 获取所有服务列表
func (nm *NacosManager) GetAllServices() (model.ServiceList, error) {

	//1.获取所有服务列表
	services, err := nm.client.GetAllServicesInfo(vo.GetAllServiceInfoParam{
		GroupName: nm.config.Nacos.Group,
		PageNo:    1,
		PageSize:  10000,
	})

	//2.如果获取失败，返回异常
	if err != nil {
		return model.ServiceList{}, err
	}

	//3.否则返回服务列表
	return services, nil
}

// GetServiceInstances 获取服务实例列表
func (nm *NacosManager) GetServiceInstances(serviceName string) ([]model.Instance, error) {

	//1.获取指定服务列表
	instances, err := nm.client.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceName,
		GroupName:   nm.config.Nacos.Group,
		HealthyOnly: true,
	})

	//2.如果获取失败，返回异常
	if err != nil {
		return nil, err
	}

	//3.否则返回服务实例列表
	return instances, nil
}
