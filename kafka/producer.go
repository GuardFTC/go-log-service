// Package kafka @Author:冯铁城 [17615007230@163.com] 2026-01-06 10:23:37
package kafka

import (
	"context"
	"errors"
	"logging-mon-service/config"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// GlobalProducer 全局生产者
var GlobalProducer *Producer

// InitProducer 初始化生产者
func InitProducer(c *config.Config) {
	GlobalProducer = NewProducer(c, context.Background())
}

// CloseProducer 关闭生产者
func CloseProducer() {
	if err := GlobalProducer.Close(); err != nil {
		logrus.Errorf("producer close error=>%v", err)
	}
}

// Producer 生产者
type Producer struct {
	w *kafka.Writer
	c context.Context
}

// NewProducer 创建生产者
func NewProducer(config *config.Config, c context.Context) *Producer {

	//1.创建生产者
	producer := &Producer{
		w: newWriterFromConfig(config),
		c: c,
	}

	//2.日志打印
	logrus.Info("producer created success")

	//3.返回
	return producer
}

// Close 关闭生产者
func (p *Producer) Close() error {

	//1.关闭生产者
	err := p.w.Close()

	//2.如果错误不为空，则返回
	if err != nil {
		return err
	}

	//3.否则打印日志，返回成功
	logrus.Info("producer closed success")
	return nil
}

// SendMassage 发送消息
func (p *Producer) SendMassage(topic string, message string) error {
	return sendMessage(topic, -1, "", message, p.w, p.c)
}

// SendMessageWithKey 发送消息（指定key）
func (p *Producer) SendMessageWithKey(topic string, key string, message string) error {
	return sendMessage(topic, -1, key, message, p.w, p.c)
}

// SendMessageWithPartition 发送消息（指定分区）
func (p *Producer) SendMessageWithPartition(topic string, partition int, message string) error {
	return sendMessage(topic, partition, "", message, p.w, p.c)
}

// SendMessages 批量发送消息
func (p *Producer) SendMessages(topic string, messages []string) error {
	return sendMessageBatch(topic, -1, "", messages, p.w, p.c)
}

// SendMessagesWithKey 批量发送消息（指定key）
func (p *Producer) SendMessagesWithKey(topic string, key string, messages []string) error {
	return sendMessageBatch(topic, -1, key, messages, p.w, p.c)
}

// SendMessagesWithPartition 批量发送消息（指定分区）
func (p *Producer) SendMessagesWithPartition(topic string, partition int, messages []string) error {
	return sendMessageBatch(topic, partition, "", messages, p.w, p.c)
}

// sendMessage 发送单条消息
func sendMessage(topic string, partition int, key string, message string, w *kafka.Writer, c context.Context) error {

	//1.校验
	if c == nil {
		return errors.New("context can not be nil")
	}

	//2.创建消息
	msg := getMessage(topic, partition, key, message)

	//3.发送消息
	if err := w.WriteMessages(c, msg); err != nil {
		return err
	}

	//4.打印日志
	logrus.Infof("producer send message=>[topic=%s partition=%d key=%s] success", topic, partition, key)

	//5.默认返回
	return nil
}

// sendMessageBatch 批量发送消息
func sendMessageBatch(topic string, partition int, key string, messages []string, w *kafka.Writer, c context.Context) error {

	//1.校验
	if c == nil {
		return errors.New("context can not be nil")
	}

	//2.创建消息切片
	var msgs []kafka.Message

	//3.循环封装消息
	for _, message := range messages {

		//4.创建消息
		msg := getMessage(topic, partition, key, message)

		//5.写入切片
		msgs = append(msgs, msg)
	}

	//6.批量发送消息
	if err := w.WriteMessages(c, msgs...); err != nil {
		return err
	}

	//7.打印日志
	logrus.Infof("producer sent %d messages=>[topic=%s partition=%d key=%s] success", len(msgs), topic, partition, key)

	//8.默认返回
	return nil
}

// getMessage 创建消息
func getMessage(topic string, partition int, key string, message string) kafka.Message {

	//1.创建消息
	msg := kafka.Message{
		Topic: topic,
		Value: []byte(message),
	}

	//2.如果key不为空，写入key
	if key != "" {
		msg.Key = []byte(key)
	}

	//3.如果分区不为空，指定分区
	if partition != -1 {
		msg.Partition = partition
	}

	//4.返回消息
	return msg
}
