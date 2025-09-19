package node

import (
	"log"
	"smart-weaver/internal/domain/agent/model/entity"
	"smart-weaver/internal/domain/agent/model/valobj"
	"smart-weaver/internal/domain/agent/service/armory"
	"smart-weaver/internal/domain/agent/service/armory/factory/context"
	"sync"
)

// Repository 接口定义
type Repository interface {
	QueryAiClientModelVOListByClientIds(clientIds []int64) ([]valobj.AiClientModelVO, error)
	QueryAiClientToolMcpVOListByClientIds(clientIds []int64) ([]valobj.AiClientToolMcpVO, error)
}

// RootNode 根节点
type RootNode struct {
	*armory.AbstractArmorySupport
	aiClientModelNode StrategyHandler
	repository        Repository
}

// NewRootNode 创建根节点实例
func NewRootNode(threadPoolSize int, repo Repository, aiClientModelNode StrategyHandler) *RootNode {
	return &RootNode{
		AbstractArmorySupport: &armory.AbstractArmorySupport{
			ThreadPool: make(chan func(), threadPoolSize),
			Deps:       make(map[string]any),
			Mu:         sync.Mutex{},
		},
		repository:        repo,
		aiClientModelNode: aiClientModelNode,
	}
}

// MultiThread 覆盖父类的多线程方法
func (r *RootNode) MultiThread(req any, ctx any) error {
	requestParameter, ok := req.(*entity.AiAgentEngineStarterEntity)
	if !ok {
		log.Println("Error: invalid request parameter type")
		return nil
	}

	dynamicContext, ok := ctx.(*context.DynamicContext)
	if !ok {
		log.Println("Error: invalid dynamic context type")
		return nil
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	// 存储结果的变量
	var aiClientModelList []valobj.AiClientModelVO
	var aiClientToolMcpList []valobj.AiClientToolMcpVO
	var err1, err2 error

	// 异步查询 ai_client_model 数据
	wg.Add(1)
	r.SubmitTask(func() {
		defer wg.Done()
		log.Printf("查询配置数据(ai_client_model) %v", requestParameter.ClientIDList)
		list, err := r.repository.QueryAiClientModelVOListByClientIds(requestParameter.ClientIDList)

		mu.Lock()
		aiClientModelList = list
		err1 = err
		mu.Unlock()
	})

	// 异步查询 ai_client_tool_mcp 数据
	wg.Add(1)
	r.SubmitTask(func() {
		defer wg.Done()
		log.Printf("查询配置数据(ai_client_tool_mcp) %v", requestParameter.ClientIDList)
		list, err := r.repository.QueryAiClientToolMcpVOListByClientIds(requestParameter.ClientIDList)

		mu.Lock()
		aiClientToolMcpList = list
		err2 = err
		mu.Unlock()
	})

	// 等待所有任务完成
	wg.Wait()

	// 检查错误
	if err1 != nil {
		log.Printf("Error querying ai_client_model: %v", err1)
		return err1
	}
	if err2 != nil {
		log.Printf("Error querying ai_client_tool_mcp: %v", err2)
		return err2
	}

	// 设置结果到动态上下文
	dynamicContext.SetValue("aiClientModelList", aiClientModelList)
	dynamicContext.SetValue("aiClientToolMcpList", aiClientToolMcpList)

	return nil
}

// DoApply 执行应用逻辑
func (r *RootNode) DoApply(requestParameter *entity.AiAgentEngineStarterEntity, dynamicContext *context.DynamicContext) (string, error) {
	log.Println("RootNode 开始执行")

	// 先执行多线程数据查询
	err := r.MultiThread(requestParameter, dynamicContext)
	if err != nil {
		log.Printf("多线程查询数据失败: %v", err)
		return "", err
	}

	// 然后路由到下一个节点
	return r.Router(requestParameter, dynamicContext)
}

// Get 获取下一个策略处理器
func (r *RootNode) Get(requestParameter *entity.AiAgentEngineStarterEntity, dynamicContext *context.DynamicContext) (StrategyHandler, error) {
	return r.aiClientModelNode, nil
}

// Router 路由方法
func (r *RootNode) Router(requestParameter *entity.AiAgentEngineStarterEntity, dynamicContext *context.DynamicContext) (string, error) {
	nextHandler, err := r.Get(requestParameter, dynamicContext)
	if err != nil {
		log.Printf("获取下一个处理器失败: %v", err)
		return "", err
	}

	if nextHandler != nil {
		log.Println("RootNode 路由到下一个处理器")
		return nextHandler.DoApply(requestParameter, dynamicContext)
	}

	// 如果没有下一个处理器，返回成功
	log.Println("RootNode 处理完成")
	return "success", nil
}

// SetAiClientModelNode 设置下一个节点（用于依赖注入）
func (r *RootNode) SetAiClientModelNode(node StrategyHandler) {
	r.aiClientModelNode = node
}

// GetRepository 获取仓库实例
func (r *RootNode) GetRepository() Repository {
	return r.repository
}
