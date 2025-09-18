package config

import (
	"fmt"
	"sync"
)

// Task 定义任务类型
type Task func()

// AsyncExecutor 异步执行器
type AsyncExecutor struct {
	corePoolSize     int
	maxPoolSize      int
	queueCapacity    int
	threadNamePrefix string
	taskQueue        chan Task
	workers          []chan Task
	once             sync.Once
}

// AsyncConfiguration 异步配置
type AsyncConfiguration struct{}

// GetAsyncExecutor 获取异步执行器
func (ac *AsyncConfiguration) GetAsyncExecutor() *AsyncExecutor {
	executor := &AsyncExecutor{
		corePoolSize:     10,
		maxPoolSize:      50,
		queueCapacity:    100,
		threadNamePrefix: "Async-Executor-",
	}
	executor.initialize()
	return executor
}

// initialize 初始化执行器
func (e *AsyncExecutor) initialize() {
	e.once.Do(func() {
		// 创建任务队列
		e.taskQueue = make(chan Task, e.queueCapacity)

		// 创建核心工作协程
		for i := 0; i < e.corePoolSize; i++ {
			go e.worker(fmt.Sprintf("%s%d", e.threadNamePrefix, i))
		}
	})
}

// worker 工作协程
func (e *AsyncExecutor) worker(name string) {
	for task := range e.taskQueue {
		task()
	}
}

// Execute 执行异步任务
func (e *AsyncExecutor) Execute(task Task) {
	e.taskQueue <- task
}
