// Package feign 负载均衡器
package feign

import (
	"math/rand"
	"sync"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/model"
)

// LoadBalancer 负载均衡器接口
type LoadBalancer interface {
	Select(instances []model.Instance) model.Instance // 选择一个实例
}

// LoadBalancerType 负载均衡类型
type LoadBalancerType int

// LoadBalancerType 负载均衡类型枚举
const (
	Random         LoadBalancerType = 0 //随机负载均衡
	RoundRobin     LoadBalancerType = 1 //轮询负载均衡
	WeightedRandom LoadBalancerType = 2 //加权随机负载均衡
)

// GetLoadBalancer 获取负载均衡器
func GetLoadBalancer(lbType LoadBalancerType) LoadBalancer {
	switch lbType {
	case Random:
		return NewRandomLoadBalancer()
	case RoundRobin:
		return NewRoundRobinLoadBalancer()
	case WeightedRandom:
		return NewWeightedRandomLoadBalancer()
	default:
		return NewRoundRobinLoadBalancer()
	}
}

//-------------------------------------------随机负载均衡器-------------------------------------------//

// RandomLoadBalancer 随机负载均衡器
type RandomLoadBalancer struct{}

// NewRandomLoadBalancer 创建随机负载均衡器
func NewRandomLoadBalancer() *RandomLoadBalancer {
	return &RandomLoadBalancer{}
}

// Select 选择一个实例
func (r *RandomLoadBalancer) Select(instances []model.Instance) model.Instance {

	//1.如果实例列表为空，返回空实例
	if len(instances) == 0 {
		return model.Instance{}
	}

	//2.如果实例列表只有一个，直接返回
	if len(instances) == 1 {
		return instances[0]
	}

	//3.从实例列表中随机选择一个实例,返回
	ran := rand.New(rand.NewSource(time.Now().UnixNano()))
	return instances[ran.Intn(len(instances))]
}

//-------------------------------------------轮询负载均衡器-------------------------------------------//

// RoundRobinLoadBalancer 轮询负载均衡器
type RoundRobinLoadBalancer struct {
	counter map[string]int
	mutex   sync.Mutex
}

// NewRoundRobinLoadBalancer 创建轮询负载均衡器
func NewRoundRobinLoadBalancer() *RoundRobinLoadBalancer {
	return &RoundRobinLoadBalancer{
		counter: make(map[string]int),
	}
}

// Select 选择一个实例
func (rr *RoundRobinLoadBalancer) Select(instances []model.Instance) model.Instance {

	//1.如果实例列表为空，返回空实例
	if len(instances) == 0 {
		return model.Instance{}
	}

	//2.如果实例列表只有一个，直接返回
	if len(instances) == 1 {
		return instances[0]
	}

	//3.加锁
	rr.mutex.Lock()
	defer rr.mutex.Unlock()

	//4.采用服务名作为key，从而实现不同服务之间策略的隔离
	key := instances[0].ServiceName

	//5.基于取模算法进行轮询
	count := rr.counter[key]
	rr.counter[key] = (count + 1) % len(instances)

	//6.获取当前轮询到的实例
	return instances[count]
}

//-------------------------------------------随机权重负载均衡器-------------------------------------------//

// WeightedRandomLoadBalancer 加权随机负载均衡器
type WeightedRandomLoadBalancer struct{}

// NewWeightedRandomLoadBalancer 创建加权随机负载均衡器
func NewWeightedRandomLoadBalancer() *WeightedRandomLoadBalancer {
	return &WeightedRandomLoadBalancer{}
}

// Select 选择一个实例
func (w *WeightedRandomLoadBalancer) Select(instances []model.Instance) model.Instance {

	//1.如果实例列表为空，返回空实例
	if len(instances) == 0 {
		return model.Instance{}
	}

	//2.如果实例列表只有一个，直接返回
	if len(instances) == 1 {
		return instances[0]
	}

	//3.计算总权重
	totalWeight := 0.0
	for _, instance := range instances {
		totalWeight += instance.Weight
	}

	//4.如果没有权重信息，使用随机策略
	ran := rand.New(rand.NewSource(time.Now().UnixNano()))
	if totalWeight <= 0 {
		return instances[ran.Intn(len(instances))]
	}

	//5.生成随机数
	random := ran.Float64() * totalWeight

	//6.根据权重选择实例
	currentWeight := 0.0
	for _, instance := range instances {
		currentWeight += instance.Weight
		if random <= currentWeight {
			return instance
		}
	}

	//7.兜底返回第一个实例
	return instances[0]
}
