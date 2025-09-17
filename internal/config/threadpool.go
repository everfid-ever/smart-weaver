package config

import (
	"sync"
)

// ThreadPool 线程池
type ThreadPool struct {
	workers     int
	workerQueue chan chan func()
	jobQueue    chan func()
	quit        chan bool
	wg          sync.WaitGroup
}

// NewThreadPool 创建线程池
func NewThreadPool(workers, queueSize int) *ThreadPool {
	return &ThreadPool{
		workers:     workers,
		workerQueue: make(chan chan func(), workers),
		jobQueue:    make(chan func(), queueSize),
		quit:        make(chan bool),
	}
}

// Start 启动线程池
func (tp *ThreadPool) Start() {
	for i := 0; i < tp.workers; i++ {
		worker := NewWorker(tp.workerQueue, tp.quit)
		worker.Start()
	}

	go tp.dispatch()
}

// Submit 提交任务
func (tp *ThreadPool) Submit(job func()) {
	tp.jobQueue <- job
}

// Stop 停止线程池
func (tp *ThreadPool) Stop() {
	close(tp.quit)
	tp.wg.Wait()
}

func (tp *ThreadPool) dispatch() {
	for {
		select {
		case job := <-tp.jobQueue:
			worker := <-tp.workerQueue
			worker <- job
		case <-tp.quit:
			return
		}
	}
}

// Worker 工作者
type Worker struct {
	workerPool chan chan func()
	jobChannel chan func()
	quit       chan bool
}

// NewWorker 创建工作者
func NewWorker(workerPool chan chan func(), quit chan bool) *Worker {
	return &Worker{
		workerPool: workerPool,
		jobChannel: make(chan func()),
		quit:       quit,
	}
}

// Start 启动工作者
func (w *Worker) Start() {
	go func() {
		for {
			w.workerPool <- w.jobChannel

			select {
			case job := <-w.jobChannel:
				job()
			case <-w.quit:
				return
			}
		}
	}()
}

// InitThreadPool 初始化线程池
func InitThreadPool(cfg *Config) *ThreadPool {
	coreSize := cfg.ThreadPool.Pool.Executor.Config.CorePoolSize
	if coreSize == 0 {
		coreSize = 20
	}

	queueSize := cfg.ThreadPool.Pool.Executor.Config.BlockQueueSize
	if queueSize == 0 {
		queueSize = 5000
	}

	tp := NewThreadPool(coreSize, queueSize)
	tp.Start()

	return tp
}
