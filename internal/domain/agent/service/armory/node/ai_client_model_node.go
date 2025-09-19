package node

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"smart-weaver/internal/domain/agent/model/entity"
	"smart-weaver/internal/domain/agent/model/valobj"
	"smart-weaver/internal/domain/agent/service/armory"
	"smart-weaver/internal/domain/agent/service/armory/factory/context"
)

// OpenAiApi OpenAI API配置（模拟Java中的OpenAiApi）
type OpenAiApi struct {
	BaseURL         string
	APIKey          string
	CompletionsPath string
	EmbeddingsPath  string
}

// OpenAiChatOptions OpenAI聊天选项（模拟Java中的OpenAiChatOptions）
type OpenAiChatOptions struct {
	Model         string
	ToolCallbacks interface{} // 工具回调
}

// OpenAiChatModel OpenAI聊天模型（模拟Java中的OpenAiChatModel）
type OpenAiChatModel struct {
	OpenAiApi      *OpenAiApi
	DefaultOptions *OpenAiChatOptions
}

// SyncMcpToolCallbackProvider 同步MCP工具回调提供者（模拟Java中的SyncMcpToolCallbackProvider）
type SyncMcpToolCallbackProvider struct {
	McpSyncClients []McpSyncClient
}

// GetToolCallbacks 获取工具回调
func (provider *SyncMcpToolCallbackProvider) GetToolCallbacks() interface{} {
	// 模拟返回工具回调集合
	callbacks := make(map[string]interface{})
	for i, client := range provider.McpSyncClients {
		callbacks[fmt.Sprintf("callback_%d", i)] = client
	}
	return callbacks
}

// NewSyncMcpToolCallbackProvider 创建同步MCP工具回调提供者
func NewSyncMcpToolCallbackProvider(mcpSyncClients []McpSyncClient) *SyncMcpToolCallbackProvider {
	return &SyncMcpToolCallbackProvider{
		McpSyncClients: mcpSyncClients,
	}
}

// OpenAiApiBuilder OpenAI API构建器
type OpenAiApiBuilder struct {
	baseURL         string
	apiKey          string
	completionsPath string
	embeddingsPath  string
}

// NewOpenAiApiBuilder 创建OpenAI API构建器
func NewOpenAiApiBuilder() *OpenAiApiBuilder {
	return &OpenAiApiBuilder{}
}

// BaseURL 设置BaseURL
func (b *OpenAiApiBuilder) BaseURL(baseURL string) *OpenAiApiBuilder {
	b.baseURL = baseURL
	return b
}

// APIKey 设置APIKey
func (b *OpenAiApiBuilder) APIKey(apiKey string) *OpenAiApiBuilder {
	b.apiKey = apiKey
	return b
}

// CompletionsPath 设置CompletionsPath
func (b *OpenAiApiBuilder) CompletionsPath(completionsPath string) *OpenAiApiBuilder {
	b.completionsPath = completionsPath
	return b
}

// EmbeddingsPath 设置EmbeddingsPath
func (b *OpenAiApiBuilder) EmbeddingsPath(embeddingsPath string) *OpenAiApiBuilder {
	b.embeddingsPath = embeddingsPath
	return b
}

// Build 构建OpenAiApi
func (b *OpenAiApiBuilder) Build() *OpenAiApi {
	return &OpenAiApi{
		BaseURL:         b.baseURL,
		APIKey:          b.apiKey,
		CompletionsPath: b.completionsPath,
		EmbeddingsPath:  b.embeddingsPath,
	}
}

// OpenAiChatOptionsBuilder OpenAI聊天选项构建器
type OpenAiChatOptionsBuilder struct {
	model         string
	toolCallbacks interface{}
}

// NewOpenAiChatOptionsBuilder 创建OpenAI聊天选项构建器
func NewOpenAiChatOptionsBuilder() *OpenAiChatOptionsBuilder {
	return &OpenAiChatOptionsBuilder{}
}

// Model 设置模型
func (b *OpenAiChatOptionsBuilder) Model(model string) *OpenAiChatOptionsBuilder {
	b.model = model
	return b
}

// ToolCallbacks 设置工具回调
func (b *OpenAiChatOptionsBuilder) ToolCallbacks(toolCallbacks interface{}) *OpenAiChatOptionsBuilder {
	b.toolCallbacks = toolCallbacks
	return b
}

// Build 构建OpenAiChatOptions
func (b *OpenAiChatOptionsBuilder) Build() *OpenAiChatOptions {
	return &OpenAiChatOptions{
		Model:         b.model,
		ToolCallbacks: b.toolCallbacks,
	}
}

// OpenAiChatModelBuilder OpenAI聊天模型构建器
type OpenAiChatModelBuilder struct {
	openAiApi      *OpenAiApi
	defaultOptions *OpenAiChatOptions
}

// NewOpenAiChatModelBuilder 创建OpenAI聊天模型构建器
func NewOpenAiChatModelBuilder() *OpenAiChatModelBuilder {
	return &OpenAiChatModelBuilder{}
}

// OpenAiApi 设置OpenAI API
func (b *OpenAiChatModelBuilder) OpenAiApi(openAiApi *OpenAiApi) *OpenAiChatModelBuilder {
	b.openAiApi = openAiApi
	return b
}

// DefaultOptions 设置默认选项
func (b *OpenAiChatModelBuilder) DefaultOptions(defaultOptions *OpenAiChatOptions) *OpenAiChatModelBuilder {
	b.defaultOptions = defaultOptions
	return b
}

// Build 构建OpenAiChatModel
func (b *OpenAiChatModelBuilder) Build() *OpenAiChatModel {
	return &OpenAiChatModel{
		OpenAiApi:      b.openAiApi,
		DefaultOptions: b.defaultOptions,
	}
}

// AiClientModelNode AI客户端模型节点
type AiClientModelNode struct {
	*armory.AbstractArmorySupport
	AiClientNode StrategyHandler
}

// NewAiClientModelNode 创建AiClientModelNode实例
func NewAiClientModelNode(aiClientNode StrategyHandler) *AiClientModelNode {
	return &AiClientModelNode{
		AbstractArmorySupport: &armory.AbstractArmorySupport{
			ThreadPool: make(chan func(), 100),
			Deps:       make(map[string]any),
		},
		AiClientNode: aiClientNode,
	}
}

// DoApply 执行应用逻辑
func (node *AiClientModelNode) DoApply(requestParameter *entity.AiAgentEngineStarterEntity, dynamicContext *context.DynamicContext) (string, error) {
	reqJSON, _ := json.Marshal(requestParameter)
	log.Printf("Ai Agent 构建，客户端构建节点 %s", string(reqJSON))

	// 从动态上下文获取AI客户端模型列表
	aiClientModelListVal := dynamicContext.GetValue("aiClientModelList")
	if aiClientModelListVal == nil {
		log.Println("没有可用的AI客户端模型配置")
		return "", nil
	}

	aiClientModelList, ok := aiClientModelListVal.([]valobj.AiClientModelVO)
	if !ok || len(aiClientModelList) == 0 {
		log.Println("没有可用的AI客户端模型配置")
		return "", nil
	}

	// 遍历模型列表，为每个模型创建对应的Bean
	for _, modelVO := range aiClientModelList {
		// 创建OpenAiChatModel对象
		chatModel, err := node.createOpenAiChatModel(modelVO)
		if err != nil {
			log.Printf("创建OpenAiChatModel失败: %v", err)
			continue
		}

		// 注册Bean
		beanName := node.beanName(modelVO.ID)
		node.RegisterDependency(beanName, chatModel)
	}

	return node.Router(requestParameter, dynamicContext)
}

// Get 获取下一个处理器
func (node *AiClientModelNode) Get(requestParameter *entity.AiAgentEngineStarterEntity, dynamicContext *context.DynamicContext) (StrategyHandler, error) {
	return node.AiClientNode, nil
}

// Router 路由到下一个处理器
func (node *AiClientModelNode) Router(requestParameter *entity.AiAgentEngineStarterEntity, dynamicContext *context.DynamicContext) (string, error) {
	nextHandler, err := node.Get(requestParameter, dynamicContext)
	if err != nil {
		return "", err
	}
	if nextHandler == nil {
		return "completed", nil
	}
	return nextHandler.DoApply(requestParameter, dynamicContext)
}

// beanName 生成Bean名称
func (node *AiClientModelNode) beanName(id int64) string {
	return "AiClientModel_" + strconv.FormatInt(id, 10)
}

// createOpenAiChatModel 创建OpenAiChatModel对象
func (node *AiClientModelNode) createOpenAiChatModel(modelVO valobj.AiClientModelVO) (*OpenAiChatModel, error) {
	// 构建OpenAiApi
	openAiApi := NewOpenAiApiBuilder().
		BaseURL(modelVO.BaseURL).
		APIKey(modelVO.APIKey).
		CompletionsPath(modelVO.CompletionsPath).
		EmbeddingsPath(modelVO.EmbeddingsPath).
		Build()

	// 收集MCP客户端
	var mcpSyncClients []McpSyncClient
	toolConfigs := modelVO.AIClientModelToolConfigs
	if len(toolConfigs) > 0 {
		for _, toolConfig := range toolConfigs {
			toolID := toolConfig.ToolID
			mcpBeanName := "AiClientToolMcp_" + strconv.FormatInt(toolID, 10)

			// 从依赖容器获取MCP客户端
			mcpClientInterface := node.GetDependency(mcpBeanName)
			if mcpClientInterface != nil {
				if mcpSyncClient, ok := mcpClientInterface.(McpSyncClient); ok {
					mcpSyncClients = append(mcpSyncClients, mcpSyncClient)
				} else {
					log.Printf("警告: Bean %s 不是McpSyncClient类型", mcpBeanName)
				}
			} else {
				log.Printf("警告: 未找到Bean %s", mcpBeanName)
			}
		}
	}

	// 创建工具回调提供者
	toolCallbackProvider := NewSyncMcpToolCallbackProvider(mcpSyncClients)

	// 构建OpenAiChatModel
	chatModel := NewOpenAiChatModelBuilder().
		OpenAiApi(openAiApi).
		DefaultOptions(
			NewOpenAiChatOptionsBuilder().
				Model(modelVO.ModelVersion).
				ToolCallbacks(toolCallbackProvider.GetToolCallbacks()).
				Build(),
		).
		Build()

	return chatModel, nil
}
