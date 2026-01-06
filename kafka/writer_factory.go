// Package kafka @Author:冯铁城 [17615007230@163.com] 2026-01-06 10:35:00
package kafka

import (
	"logging-mon-service/config"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// newWriterFromConfig 从配置创建Writer
func newWriterFromConfig(config *config.Config) *kafka.Writer {
	return &kafka.Writer{

		//1.常规配置
		Addr:        kafka.TCP(config.Kafka.Brokers...),
		Async:       config.Kafka.Async,
		Compression: getCompression(config.Kafka.Compression),
		Balancer:    getBalancer(config.Kafka.Balancer),

		//2.发送确认配置
		RequiredAcks: getRequiredAcks(config.Kafka.RequiredAcks),
		MaxAttempts:  config.Kafka.MaxAttempts,

		//3.超时配置
		WriteTimeout: config.Kafka.GetWriteTimeout(),
		ReadTimeout:  config.Kafka.GetReadTimeout(),

		//4.消息批处理阈值配置
		BatchSize:    config.Kafka.BatchSize,
		BatchTimeout: config.Kafka.GetBatchTimeout(),
		BatchBytes:   config.Kafka.BatchBytes,

		//5.网络配置
		Transport: &kafka.Transport{
			DialTimeout: config.Kafka.GetDialTimeout(),
			IdleTimeout: config.Kafka.GetIdleTimeout(),
		},

		//6.日志配置
		ErrorLogger: kafka.LoggerFunc(logrus.Errorf),
	}
}

// getCompression 获取压缩算法
func getCompression(compression string) kafka.Compression {
	switch compression {
	case "gzip":
		return kafka.Gzip
	case "snappy":
		return kafka.Snappy
	case "lz4":
		return kafka.Lz4
	case "zstd":
		return kafka.Zstd
	default:
		return kafka.Snappy
	}
}

// getBalancer 获取负载均衡策略
func getBalancer(balancer string) kafka.Balancer {
	switch balancer {
	case "least_bytes":
		return &kafka.LeastBytes{}
	case "round_robin":
		return &kafka.RoundRobin{}
	case "hash":
		return &kafka.Hash{}
	case "crc32":
		return &kafka.CRC32Balancer{}
	default:
		return &kafka.RoundRobin{}
	}
}

// getRequiredAcks 获取确认级别
func getRequiredAcks(acks int) kafka.RequiredAcks {
	switch acks {
	case -1:
		return kafka.RequireAll
	case 0:
		return kafka.RequireNone
	case 1:
		return kafka.RequireOne
	default:
		return kafka.RequireAll
	}
}
