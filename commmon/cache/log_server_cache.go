// Package cache @Author:冯铁城 [17615007230@163.com] 2025-12-30 14:27:29
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"logging-mon-service/model"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// LogServerCacheManager 日志服务缓存管理器
type LogServerCacheManager struct {
	cacheFile string              // 缓存文件路径
	service   *LogServerService   // 服务客户端
	cache     *model.LogServerObj // 内存缓存
	mutex     sync.RWMutex        // 读写锁
	ticker    *time.Ticker        // 定时器
	ctx       context.Context     // 上下文
	cancel    context.CancelFunc  // 取消函数
}

// NewLogServerCacheManager 创建日志服务缓存管理器
func NewLogServerCacheManager(serviceName string) *LogServerCacheManager {

	//1.创建缓存文件
	cacheFile := filepath.Join(os.TempDir(), fmt.Sprintf("LogServerObj-%s.json", serviceName))

	//2.创建取消上下文
	ctx, cancel := context.WithCancel(context.Background())

	//3.创建日志服务缓存管理器
	manager := &LogServerCacheManager{
		cacheFile: cacheFile,
		service:   NewLogServerService(),
		ctx:       ctx,
		cancel:    cancel,
	}

	//4.初始化缓存
	manager.initialize()

	//5.返回
	return manager
}

//----------------------------------------------外部方法---------------------------------------------------//

// Start 启动定时更新任务
func (m *LogServerCacheManager) Start() {

	//1.创建定时器,30s执行一次
	m.ticker = time.NewTicker(30 * time.Second)

	//2.创建协程，执行定时任务
	go func() {
		for {
			select {
			case <-m.ticker.C: //监听定时器chan，每30s执行一次
				m.updateCache()
			case <-m.ctx.Done(): //监听取消信号，如果取消，则退出
				return
			}
		}
	}()

	//3.打印日志
	log.Printf("[内存缓存-LogServer] 定时更新任务已启动，间隔30秒")
}

// Stop 停止定时更新任务
func (m *LogServerCacheManager) Stop() {

	//1.如果定时器不为nil，则停止
	if m.ticker != nil {
		m.ticker.Stop()
	}

	//2.调用上下文取消函数，让定时任务协程退出
	m.cancel()

	//3.打印日志
	log.Printf("[内存缓存-LogServer] 定时更新任务已停止")
}

// GetLogServerObj 获取日志服务对象（从内存缓存）
func (m *LogServerCacheManager) GetLogServerObj() *model.LogServerObj {

	//1.加读锁，确保线程安全
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	//2.如果内存缓存为nil，则返回nil
	if m.cache == nil {
		return nil
	}

	//3.创建副本，避免直接返回对象被外部更改
	result := model.NewLogServerObj(len(m.cache.ProjectObjs))
	copy(result.ProjectObjs, m.cache.ProjectObjs)

	//4.返回副本
	return result
}

// GetCacheFilePath 获取缓存文件路径
func (m *LogServerCacheManager) GetCacheFilePath() string {
	return m.cacheFile
}

// ForceUpdate 强制更新缓存
func (m *LogServerCacheManager) ForceUpdate() {
	m.updateCache()
}

//----------------------------------------------内部方法---------------------------------------------------//

// initialize 初始化缓存
func (m *LogServerCacheManager) initialize() {

	//1.优先从文件读取
	if obj := m.readFromFile(); obj != nil {

		//2.内存缓存写操作加锁，确保线程安全
		m.mutex.Lock()
		m.cache = obj
		m.mutex.Unlock()

		//3.打印日志，返回
		log.Printf("[内存缓存-LogServer] 从缓存文件加载成功")
		return
	}

	//4.文件读取失败，尝试从HTTP获取
	if obj, err := m.service.GetLogServerObj(); err == nil {

		//5.内存缓存写操作加锁，确保线程安全
		m.mutex.Lock()
		m.cache = obj
		m.mutex.Unlock()

		//6.将数据保存到文件
		m.saveToFile(obj)

		//7.打印日志，返回
		log.Printf("[内存缓存-LogServer] 从HTTP接口初始化成功")
	} else {
		log.Printf("[内存缓存-LogServer] 初始化失败:[%v]", err)
	}
}

// readFromFile 从文件读取缓存
func (m *LogServerCacheManager) readFromFile() *model.LogServerObj {

	//1.检查文件是否存在
	if _, err := os.Stat(m.cacheFile); os.IsNotExist(err) {
		return nil
	}

	//2.读取文件内容
	data, err := os.ReadFile(m.cacheFile)
	if err != nil {
		log.Printf("[内存缓存-LogServer] 读取缓存文件失败: %v", err)
		return nil
	}

	//3.JSON反序列化
	var obj model.LogServerObj
	if err := json.Unmarshal(data, &obj); err != nil {
		log.Printf("[内存缓存-LogServer] 解析缓存文件失败: %v", err)
		return nil
	}

	//4.返回
	return &obj
}

// saveToFile 保存到文件
func (m *LogServerCacheManager) saveToFile(obj *model.LogServerObj) {

	//1.如果对象为nil，则返回
	if obj == nil {
		return
	}

	//2.JSON序列化
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		log.Printf("[内存缓存-LogServer] 序列化对象失败: %v", err)
		return
	}

	//3.原子写入:先写临时文件
	tempFile := m.cacheFile + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		log.Printf("[内存缓存-LogServer] 写入临时文件失败: %v", err)
		return
	}

	//4.重命名为正式文件，如果失败，则清理临时文件
	if err := os.Rename(tempFile, m.cacheFile); err != nil {
		log.Printf("[内存缓存-LogServer] 重命名缓存文件失败: %v", err)
		os.Remove(tempFile)
		return
	}

	//5.打印日志
	log.Printf("[内存缓存-LogServer] 缓存文件保存成功: %s", m.cacheFile)
}

// updateCache 更新缓存
func (m *LogServerCacheManager) updateCache() {

	//1.调用HTTP接口获取最新数据
	obj, err := m.service.GetLogServerObj()
	if err != nil {
		log.Printf("[内存缓存-LogServer] 更新缓存失败: %v", err)
		return
	}

	//2.更新内存缓存，写操作加锁，确保线程安全
	m.mutex.Lock()
	m.cache = obj
	m.mutex.Unlock()

	//3.保存到文件
	m.saveToFile(obj)

	//4.日志打印
	log.Printf("[内存缓存-LogServer] 缓存更新成功")
}
