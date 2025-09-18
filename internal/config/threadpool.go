package config

import (
	"log"
	"sync"
	"time"
)

type ThreadPoolConfigProperties struct {
	CorePoolSize   int    `yaml:"corePoolSize"`
	MaxPoolSize    int    `yaml:"maxPoolSize"`
	KeepAliveTime  int64  `yaml:"keepAliveTime"` // 秒
	BlockQueueSize int    `yaml:"blockQueueSize"`
	Policy         string `yaml:"policy"` // AbortPolicy, DiscardPolicy, DiscardOldestPolicy, CallerRunsPolicy
}

// ThreadPoolExecutor 模拟 Java 的 ThreadPoolExecutor
type ThreadPoolExecutor struct {
	corePoolSize  int
	maxPoolSize   int
	keepAliveTime time.Duration
	queue         chan Task
	policy        string

	wg       sync.WaitGroup
	shutdown chan struct{}
}

func NewThreadPoolExecutor(properties ThreadPoolConfigProperties) *ThreadPoolExecutor {
	pool := &ThreadPoolExecutor{
		corePoolSize:  properties.CorePoolSize,
		maxPoolSize:   properties.MaxPoolSize,
		keepAliveTime: time.Duration(properties.KeepAliveTime) * time.Second,
		queue:         make(chan Task, properties.BlockQueueSize),
		policy:        properties.Policy,
		shutdown:      make(chan struct{}),
	}

	// 启动核心线程
	for i := 0; i < pool.corePoolSize; i++ {
		go pool.worker()
	}
	return pool
}

// worker 执行任务
func (p *ThreadPoolExecutor) worker() {
	for {
		select {
		case task := <-p.queue:
			if task != nil {
				task()
				p.wg.Done()
			}
		case <-p.shutdown:
			return
		}
	}
}

// Submit 提交任务
func (p *ThreadPoolExecutor) Submit(task Task) {
	p.wg.Add(1)
	select {
	case p.queue <- task:
		// 成功入队
	default:
		switch p.policy {
		case "AbortPolicy":
			log.Panic("Task rejected: AbortPolicy")
		case "DiscardPolicy":
			p.wg.Done()
		case "DiscardOldestPolicy":
			select {
			case <-p.queue:
				p.queue <- task
			default:
				p.wg.Done()
			}
		case "CallerRunsPolicy":
			task()
			p.wg.Done()
		default:
			log.Panic("Task rejected: default AbortPolicy")
		}
	}
}

// Shutdown 等待任务完成并关闭
func (p *ThreadPoolExecutor) Shutdown() {
	p.wg.Wait()
	close(p.shutdown)
}

// InitThreadPool 根据配置初始化线程池
func InitThreadPool(cfg *Config) *ThreadPoolExecutor {
	execDetail := cfg.ThreadPool.Pool.Executor.Config

	properties := ThreadPoolConfigProperties{
		CorePoolSize:   execDetail.CorePoolSize,
		MaxPoolSize:    execDetail.MaxPoolSize,
		KeepAliveTime:  execDetail.KeepAliveTime,
		BlockQueueSize: execDetail.BlockQueueSize,
		Policy:         execDetail.Policy,
	}

	return NewThreadPoolExecutor(properties)
}
