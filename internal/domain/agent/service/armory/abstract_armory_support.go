package armory

import (
	"log"
	"sync"
)

// AbstractArmorySupport 抽象生成器类
type AbstractArmorySupport struct {
	ThreadPool chan func()
	Deps       map[string]any
	Mu         sync.Mutex
}

// RegisterDependency 注册依赖对象（线程安全），已存在则覆盖
func (a *AbstractArmorySupport) RegisterDependency(name string, instance any) {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	a.Deps[name] = instance
	log.Printf("成功注册依赖: %s", name)
}

// GetDependency 获取依赖对象（线程安全），不存在返回 nil
func (a *AbstractArmorySupport) GetDependency(name string) any {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	return a.Deps[name]
}

// SubmitTask 提交任务到线程池
func (a *AbstractArmorySupport) SubmitTask(task func()) {
	a.ThreadPool <- task
}

// CloseThreadPool 关闭线程池
func (a *AbstractArmorySupport) CloseThreadPool() {
	close(a.ThreadPool)
}
