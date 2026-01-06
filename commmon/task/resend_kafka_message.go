// Package task @Author:冯铁城 [17615007230@163.com] 2026-01-06 15:17:06
package task

import (
	"logging-mon-service/commmon/util/message"
	"logging-mon-service/config"
	"logging-mon-service/kafka"
	"os"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

// resendKafkaMessageTask 重新发送Kafka消息任务
type resendKafkaMessageTask struct {
	cfg  *config.Config
	cron *cron.Cron
}

// newResendKafkaMessageTask 创建重发Kafka消息任务
func newResendKafkaMessageTask(cfg *config.Config) *resendKafkaMessageTask {

	//1.创建任务
	task := &resendKafkaMessageTask{
		cfg:  cfg,
		cron: cron.New(cron.WithSeconds()), // 启用秒级 cron
	}

	//2.添加定时任务：每小时第1分钟（0 1 * * * *）
	//执行时间示例
	//2026-01-05 15:01:00
	//2026-01-05 16:01:00
	_, err := task.cron.AddFunc("0 1 * * * *", task.resendKafkaMessage)
	if err != nil {
		logrus.Errorf("[定时任务] 消息补发 创建失败: %v", err)
	}

	//3.返回任务实例
	return task
}

// Start 启动定时任务
func (t *resendKafkaMessageTask) start() {
	t.cron.Start()
	logrus.Infof("[定时任务] 消息补发 已启动，每小时01分执行")
}

// Stop 停止定时任务
func (t *resendKafkaMessageTask) stop() {
	t.cron.Stop()
	logrus.Infof("[定时任务] 消息补发 已停止")
}

// resendKafkaMessage 核心重发逻辑
func (t *resendKafkaMessageTask) resendKafkaMessage() {

	//1.获取上一个小时 格式:yyyyMMddHH
	hour := time.Now().Add(-1 * time.Hour).Format("2006010215")
	logrus.Infof("[消息补发] 开始 时间范围[%s]", hour)

	//2.读取失败消息文件
	files, err := message.ReadMessageFiles(hour)
	if err != nil {
		logrus.Errorf("[消息补发] 读取消息文件失败:[%v]", err)
		return
	}
	logrus.Infof("[消息补发] 读取失败消息文件数量: [%d]", len(files))

	//3.从文件读取消息内容
	messages := t.readMessageFromFile(files)
	logrus.Infof("[消息补发] 从文件读取消息数量: [%d]", len(messages))

	//4.删除整个小时目录（包括所有文件）
	if err = message.DeleteMessageFiles(hour); err != nil {
		logrus.Errorf("[消息补发] 删除失败消息文件失败: %v", err)
	}

	//5.如果消息数量为0，则返回
	if len(messages) == 0 {
		return
	}

	//6.循环消息
	for _, resendMessage := range messages {

		//7.发送消息
		//发送失败，写入文件
		//发送成功，打印日志
		if err = kafka.GlobalProducer.SendMassage(t.cfg.Kafka.Topic, resendMessage); err != nil {
			if file, e := message.WriteMessageFile(resendMessage); e != nil {
				logrus.Errorf("[消息补发] 发送消息:[%v] 失败:[%v] 写入文件失败:[%v]", resendMessage, err, e)
			} else {
				logrus.Errorf("[消息补发] 发送消息:[%v] 失败:[%v] 写入文件成功:[%v]", resendMessage, err, file)
			}
		} else {
			logrus.Infof("[消息补发] 发送消息:[%v]成功", resendMessage)
		}
	}

	//8.打印结束日志
	logrus.Infof("[消息补发] 结束 时间范围[%s]", hour)
}

// readMessageFromFile 从文件列表读取消息内容
func (t *resendKafkaMessageTask) readMessageFromFile(files []string) []string {

	//1.定义消息切片
	var messages []string

	//2.遍历文件
	for _, path := range files {

		//3.读取文件内容
		content, err := os.ReadFile(path)
		if err != nil {
			logrus.Errorf("[消息补发] 读取文件失败 [%s]: %v", path, err)
			continue
		}

		//4.解析消息为字符串
		resendMessage := string(content)
		if resendMessage == "" {
			continue
		}

		//5.存入切片
		messages = append(messages, resendMessage)
	}

	//6.返回消息切片
	return messages
}
