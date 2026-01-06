// Package config @Author:冯铁城 [17615007230@163.com] 2026-01-06 10:30:00
package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// KafkaConfig Kafka配置
type KafkaConfig struct {

	//基础配置
	Brokers     []string `json:"brokers"`     // Broker地址列表
	Async       bool     `json:"async"`       // 是否异步发送
	Compression string   `json:"compression"` // 压缩算法: gzip, snappy, lz4, zstd
	Balancer    string   `json:"balancer"`    // 负载均衡策略: round_robin, least_bytes, hash, crc32, consistent_hash

	//发送确认配置
	RequiredAcks int `json:"required_acks"` // 确认级别: -1(all), 0(none), 1(leader)
	MaxAttempts  int `json:"max_attempts"`  // 最大重试次数

	//超时配置
	WriteTimeout int `json:"write_timeout"` // 写超时(秒)
	ReadTimeout  int `json:"read_timeout"`  // 读超时(秒)

	//批处理配置
	BatchSize    int   `json:"batch_size"`    // 批处理消息数量
	BatchTimeout int   `json:"batch_timeout"` // 批处理超时(毫秒)
	BatchBytes   int64 `json:"batch_bytes"`   // 批处理字节数

	//网络配置
	DialTimeout int `json:"dial_timeout"` // 连接超时(秒)
	IdleTimeout int `json:"idle_timeout"` // 空闲超时(秒)
}

// parseKafkaConfig 解析Kafka配置
func parseKafkaConfig(c *Config) {

	//1.解析Brokers
	if envBrokers := os.Getenv("KAFKA_BROKERS"); envBrokers != "" {
		c.Kafka.Brokers = strings.Split(envBrokers, ",")
	}

	//2.解析异步发送
	if envAsync := os.Getenv("KAFKA_ASYNC"); envAsync != "" {
		if async, err := strconv.ParseBool(envAsync); err == nil {
			c.Kafka.Async = async
		}
	}

	//3.解析压缩算法
	if envCompression := os.Getenv("KAFKA_COMPRESSION"); envCompression != "" {
		c.Kafka.Compression = envCompression
	}

	//4.解析负载均衡策略
	if envBalancer := os.Getenv("KAFKA_BALANCER"); envBalancer != "" {
		c.Kafka.Balancer = envBalancer
	}

	//5.解析确认级别
	if envAcks := os.Getenv("KAFKA_REQUIRED_ACKS"); envAcks != "" {
		if acks, err := strconv.Atoi(envAcks); err == nil {
			c.Kafka.RequiredAcks = acks
		}
	}

	//6.解析最大重试次数
	if envAttempts := os.Getenv("KAFKA_MAX_ATTEMPTS"); envAttempts != "" {
		if attempts, err := strconv.Atoi(envAttempts); err == nil {
			c.Kafka.MaxAttempts = attempts
		}
	}

	//7.解析写超时
	if envWriteTimeout := os.Getenv("KAFKA_WRITE_TIMEOUT"); envWriteTimeout != "" {
		if timeout, err := strconv.Atoi(envWriteTimeout); err == nil {
			c.Kafka.WriteTimeout = timeout
		}
	}

	//8.解析读超时
	if envReadTimeout := os.Getenv("KAFKA_READ_TIMEOUT"); envReadTimeout != "" {
		if timeout, err := strconv.Atoi(envReadTimeout); err == nil {
			c.Kafka.ReadTimeout = timeout
		}
	}

	//9.解析批处理大小
	if envBatchSize := os.Getenv("KAFKA_BATCH_SIZE"); envBatchSize != "" {
		if size, err := strconv.Atoi(envBatchSize); err == nil {
			c.Kafka.BatchSize = size
		}
	}

	//10.解析批处理超时
	if envBatchTimeout := os.Getenv("KAFKA_BATCH_TIMEOUT"); envBatchTimeout != "" {
		if timeout, err := strconv.Atoi(envBatchTimeout); err == nil {
			c.Kafka.BatchTimeout = timeout
		}
	}

	//11.解析批处理字节数
	if envBatchBytes := os.Getenv("KAFKA_BATCH_BYTES"); envBatchBytes != "" {
		if bytes, err := strconv.Atoi(envBatchBytes); err == nil {
			c.Kafka.BatchBytes = int64(bytes)
		}
	}

	//12.解析连接超时
	if envDialTimeout := os.Getenv("KAFKA_DIAL_TIMEOUT"); envDialTimeout != "" {
		if timeout, err := strconv.Atoi(envDialTimeout); err == nil {
			c.Kafka.DialTimeout = timeout
		}
	}

	//13.解析空闲超时
	if envIdleTimeout := os.Getenv("KAFKA_IDLE_TIMEOUT"); envIdleTimeout != "" {
		if timeout, err := strconv.Atoi(envIdleTimeout); err == nil {
			c.Kafka.IdleTimeout = timeout
		}
	}

	//14.设置默认值
	setKafkaDefaults(&c.Kafka)
}

// setKafkaDefaults 设置Kafka默认配置
func setKafkaDefaults(kafka *KafkaConfig) {

	//1.默认Brokers
	if len(kafka.Brokers) == 0 {
		kafka.Brokers = []string{"localhost:9092"}
	}

	//2.默认异步发送
	kafka.Async = false

	//3.默认压缩算法
	if kafka.Compression == "" {
		kafka.Compression = "snappy"
	}

	//4.默认负载均衡策略
	if kafka.Balancer == "" {
		kafka.Balancer = "round_robin"
	}

	//5.默认确认级别
	if kafka.RequiredAcks == 0 {
		kafka.RequiredAcks = -1 // RequireAll
	}

	//6.默认最大重试次数
	if kafka.MaxAttempts == 0 {
		kafka.MaxAttempts = 5
	}

	//7.默认写超时
	if kafka.WriteTimeout == 0 {
		kafka.WriteTimeout = 30
	}

	//8.默认读超时
	if kafka.ReadTimeout == 0 {
		kafka.ReadTimeout = 10
	}

	//9.默认批处理大小
	if kafka.BatchSize == 0 {
		kafka.BatchSize = 100
	}

	//10.默认批处理超时
	if kafka.BatchTimeout == 0 {
		kafka.BatchTimeout = 100
	}

	//11.默认批处理字节数
	if kafka.BatchBytes == 0 {
		kafka.BatchBytes = 1048576 // 1MB
	}

	//12.默认连接超时
	if kafka.DialTimeout == 0 {
		kafka.DialTimeout = 5
	}

	//13.默认空闲超时
	if kafka.IdleTimeout == 0 {
		kafka.IdleTimeout = 300
	}
}

// GetWriteTimeout 获取写超时时间
func (k *KafkaConfig) GetWriteTimeout() time.Duration {
	return time.Duration(k.WriteTimeout) * time.Second
}

// GetReadTimeout 获取读超时时间
func (k *KafkaConfig) GetReadTimeout() time.Duration {
	return time.Duration(k.ReadTimeout) * time.Second
}

// GetBatchTimeout 获取批处理超时时间
func (k *KafkaConfig) GetBatchTimeout() time.Duration {
	return time.Duration(k.BatchTimeout) * time.Millisecond
}

// GetDialTimeout 获取连接超时时间
func (k *KafkaConfig) GetDialTimeout() time.Duration {
	return time.Duration(k.DialTimeout) * time.Second
}

// GetIdleTimeout 获取空闲超时时间
func (k *KafkaConfig) GetIdleTimeout() time.Duration {
	return time.Duration(k.IdleTimeout) * time.Second
}
