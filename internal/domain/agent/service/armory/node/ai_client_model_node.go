package node

import (
	"fmt"
	"log"
	"smart-weaver/internal/domain/agent/model/entity"
	"smart-weaver/internal/domain/agent/model/valobj"
	"smart-weaver/internal/domain/agent/service/armory"
	"smart-weaver/internal/domain/agent/service/armory/factory/context"
)

// OpenAiChatModel OpenAI聊天模型结构体
type OpenAiChatModel struct {
	ModelName       string
	BaseURL         string
	APIKey          string
	CompletionsPath string
	EmbeddingsPath  string
	ModelType       string
	ModelVersion    string
	Timeout         int
}

// ModelRegistry 模型注册器接口
type ModelRegistry interface {
	RegisterModel(beanName string, model *OpenAiChatModel) error
	RemoveModel(beanName string) error
	ContainsModel(beanName string) bool
	GetModel(beanName string) *OpenAiChatModel
}

// DefaultModelRegistry 默认模型注册器实现
type DefaultModelRegistry struct {
	models map[string]*OpenAiChatModel
}

// NewDefaultModelRegistry 创建默认模型注册器
func NewDefaultModelRegistry() *DefaultModelRegistry {
	return &DefaultModelRegistry{
		models: make(map[string]*OpenAiChatModel),
	}
}

// RegisterModel 注册模型
func (r *DefaultModelRegistry) RegisterModel(beanName string, model *OpenAiChatModel) error {
	r.models[beanName] = model
	return nil
}

// RemoveModel 移除模型
func (r *DefaultModelRegistry) RemoveModel(beanName string) error {
	delete(r.models, beanName)
	return nil
}

// ContainsModel 检查模型是否存在
func (r *DefaultModelRegistry) ContainsModel(beanName string) bool {
	_, exists := r.models[beanName]
	return exists
}

// GetModel 获取模型
func (r *DefaultModelRegistry) GetModel(beanName string) *OpenAiChatModel {
	return r.models[beanName]
}

// AiClientModelNode AI客户端模型节点
type AiClientModelNode struct {
	*armory.AbstractArmorySupport
	aiClientToolMcpNode StrategyHandler
	modelRegistry       ModelRegistry
}

// NewAiClientModelNode 创建AI客户端模型节点
func NewAiClientModelNode(aiClientToolMcpNode StrategyHandler) *AiClientModelNode {
	return &AiClientModelNode{
		AbstractArmorySupport: &armory.AbstractArmorySupport{
			ThreadPool: make(chan func(), 100), // 创建容量为100的线程池
			Deps:       make(map[string]any),
		},
		aiClientToolMcpNode: aiClientToolMcpNode,
		modelRegistry:       NewDefaultModelRegistry(),
	}
}

// DoApply 执行应用逻辑
func (n *AiClientModelNode) DoApply(requestParameter *entity.AiAgentEngineStarterEntity, dynamicContext *context.DynamicContext) (string, error) {
	log.Println("AiAgent 装配，客户端模型")

	// 从动态上下文获取AI客户端模型列表
	aiClientModelListInterface := dynamicContext.GetValue("aiClientModelList")
	if aiClientModelListInterface == nil {
		log.Println("warn: 没有可用的AI客户端模型配置")
		// 即使没有模型配置，也继续执行下一个节点
		return n.Router(requestParameter, dynamicContext)
	}

	aiClientModelList, ok := aiClientModelListInterface.([]valobj.AiClientModelVO)
	if !ok || len(aiClientModelList) == 0 {
		log.Println("warn: 没有可用的AI客户端模型配置")
		// 即使没有模型配置，也继续执行下一个节点
		return n.Router(requestParameter, dynamicContext)
	}

	// 遍历模型列表，为每个模型创建对应的实例
	for _, modelVO := range aiClientModelList {
		// 构建Bean名称
		beanName := fmt.Sprintf("AiClientModel%d", modelVO.ID)

		// 创建OpenAiChatModel对象
		chatModel, err := n.createOpenAiChatModel(modelVO)
		if err != nil {
			log.Printf("error: 创建OpenAiChatModel失败: %v", err)
			continue
		}

		// 如果模型已存在，先移除
		if n.modelRegistry.ContainsModel(beanName) {
			if err := n.modelRegistry.RemoveModel(beanName); err != nil {
				log.Printf("error: 移除现有模型失败: %v", err)
			}
		}

		// 注册新的模型
		if err := n.modelRegistry.RegisterModel(beanName, chatModel); err != nil {
			log.Printf("error: 注册模型失败: %v", err)
			continue
		}

		log.Printf("成功注册AI客户端模型: %s", beanName)
	}

	return n.Router(requestParameter, dynamicContext)
}

// Get 获取下一个处理器
func (n *AiClientModelNode) Get(requestParameter *entity.AiAgentEngineStarterEntity, dynamicContext *context.DynamicContext) (StrategyHandler, error) {
	return n.aiClientToolMcpNode, nil
}

// Router 路由方法
func (n *AiClientModelNode) Router(requestParameter *entity.AiAgentEngineStarterEntity, dynamicContext *context.DynamicContext) (string, error) {
	nextHandler, err := n.Get(requestParameter, dynamicContext)
	if err != nil {
		return "", err
	}
	if nextHandler != nil {
		log.Println("AiClientModelNode 路由到下一个处理器")
		return nextHandler.DoApply(requestParameter, dynamicContext)
	}
	return "success", nil
}

// createOpenAiChatModel 创建OpenAI聊天模型
func (n *AiClientModelNode) createOpenAiChatModel(modelVO valobj.AiClientModelVO) (*OpenAiChatModel, error) {
	if modelVO.ModelName == "" {
		return nil, fmt.Errorf("model name cannot be empty")
	}

	if modelVO.APIKey == "" {
		return nil, fmt.Errorf("API key cannot be empty")
	}

	chatModel := &OpenAiChatModel{
		ModelName:       modelVO.ModelName,
		BaseURL:         modelVO.BaseURL,
		APIKey:          modelVO.APIKey,
		CompletionsPath: modelVO.CompletionsPath,
		EmbeddingsPath:  modelVO.EmbeddingsPath,
		ModelType:       modelVO.ModelType,
		ModelVersion:    modelVO.ModelVersion,
		Timeout:         modelVO.Timeout,
	}

	// 设置默认值
	if chatModel.BaseURL == "" {
		chatModel.BaseURL = "https://api.openai.com"
	}
	if chatModel.CompletionsPath == "" {
		chatModel.CompletionsPath = "/v1/chat/completions"
	}
	if chatModel.EmbeddingsPath == "" {
		chatModel.EmbeddingsPath = "/v1/embeddings"
	}
	if chatModel.Timeout <= 0 {
		chatModel.Timeout = 30 // 默认30秒超时
	}

	return chatModel, nil
}

// GetModelRegistry 获取模型注册器（用于测试或外部访问）
func (n *AiClientModelNode) GetModelRegistry() ModelRegistry {
	return n.modelRegistry
}

// SetAiClientToolMcpNode 设置下一个节点（用于依赖注入）
func (n *AiClientModelNode) SetAiClientToolMcpNode(node StrategyHandler) {
	n.aiClientToolMcpNode = node
}
