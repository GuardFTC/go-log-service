// Package message @Author:冯铁城 [17615007230@163.com] 2026-01-06 14:50:44
package message

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// messageDirPath 消息文件目录
var messageDirPath = filepath.Join(os.TempDir(), "logging-mon")

// WriteMessageFile 写入消息文件
func WriteMessageFile(message string) (string, error) {

	//1.获取当前小时 示例:2023010115
	hour := time.Now().Format("2006010215")

	//2.写入消息文件到指定的小时目录
	return WriteMessageFileByHour(hour, message)
}

// WriteMessageFileByHour 写入消息文件到指定的小时目录
func WriteMessageFileByHour(hour, message string) (string, error) {

	//1.拼接小时目录路径
	hourDirPath := filepath.Join(messageDirPath, hour)

	//2.如果目录不存在，创建目录（包括父目录）
	if _, err := os.Stat(hourDirPath); os.IsNotExist(err) {
		if err = os.MkdirAll(hourDirPath, 0755); err != nil {
			return "", fmt.Errorf("failed to create hour directory: %w", err)
		}
	}

	//3.生成随机文件名（UUID + .message）
	fileName := uuid.New().String() + ".message"

	//4.拼接完整文件路径
	filePath := filepath.Join(hourDirPath, fileName)

	//5.写入 UTF-8 文件
	if err := os.WriteFile(filePath, []byte(message), 0644); err != nil {
		return "", fmt.Errorf("failed to write message file: %w", err)
	}

	//6.返回文件路径
	return filePath, nil
}

// ReadMessageFiles 读取某小时目录下的所有消息文件
func ReadMessageFiles(hour string) ([]string, error) {

	//1.拼接小时目录路径
	hourDirPath := filepath.Join(messageDirPath, hour)

	//2.判断目录是否存在
	if _, err := os.Stat(hourDirPath); os.IsNotExist(err) {
		return []string{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to stat hour directory: %w", err)
	}

	//3.读取目录下所有文件
	entries, err := os.ReadDir(hourDirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read hour directory: %w", err)
	}

	//4.循环拼接文件路径
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, filepath.Join(hourDirPath, entry.Name()))
		}
	}

	//5.返回文件列表
	return files, nil
}
