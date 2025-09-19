package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"smart-weaver/internal/domain/agent/model/entity"
	"smart-weaver/internal/domain/agent/model/valobj"
	"smart-weaver/internal/domain/agent/service/armory"
	"smart-weaver/internal/domain/agent/service/armory/factory/context"
)

// McpSyncClient MCP同步客户端接口（模拟Java中的McpSyncClient）
type McpSyncClient interface {
	Initialize() (interface{}, error)
	SetRequestTimeout(timeout time.Duration)
}

// HttpClientSseClientTransport SSE传输客户端（模拟Java中的HttpClientSseClientTransport）
type HttpClientSseClientTransport struct {
	BaseURI     string
	SseEndpoint string
}

// StdioClientTransport Stdio传输客户端（模拟Java中的StdioClientTransport）
type StdioClientTransport struct {
	ServerParams *ServerParameters
}

// ServerParameters 服务器参数（模拟Java中的ServerParameters）
type ServerParameters struct {
	Command string
	Args    []string
}

// MockMcpSyncClient 模拟MCP客户端实现
type MockMcpSyncClient struct {
	transport      interface{}
	requestTimeout time.Duration
}

func (c *MockMcpSyncClient) Initialize() (interface{}, error) {
	return map[string]interface{}{"status": "initialized"}, nil
}

func (c *MockMcpSyncClient) SetRequestTimeout(timeout time.Duration) {
	c.requestTimeout = timeout
}

// AiClientToolMcpNode Tool MCP节点
type AiClientToolMcpNode struct {
	*armory.AbstractArmorySupport
	AiClientAdvisorNode StrategyHandler
}

// NewAiClientToolMcpNode 创建AiClientToolMcpNode实例
func NewAiClientToolMcpNode(aiClientAdvisorNode StrategyHandler) *AiClientToolMcpNode {
	return &AiClientToolMcpNode{
		AbstractArmorySupport: &armory.AbstractArmorySupport{
			ThreadPool: make(chan func(), 100),
			Deps:       make(map[string]any),
		},
		AiClientAdvisorNode: aiClientAdvisorNode,
	}
}

// DoApply 执行应用逻辑
func (node *AiClientToolMcpNode) DoApply(requestParameter *entity.AiAgentEngineStarterEntity, dynamicContext *context.DynamicContext) (string, error) {
	reqJSON, _ := json.Marshal(requestParameter)
	log.Printf("Ai Agent 构建，tool mcp 节点 %s", string(reqJSON))

	// 从动态上下文获取AI客户端工具MCP列表
	aiClientToolMcpListVal := dynamicContext.GetValue("aiClientToolMcpList")
	if aiClientToolMcpListVal == nil {
		log.Println("没有可用的AI客户端工具配置 MCP")
		return node.Router(requestParameter, dynamicContext)
	}

	aiClientToolMcpList, ok := aiClientToolMcpListVal.([]valobj.AiClientToolMcpVO)
	if !ok || len(aiClientToolMcpList) == 0 {
		log.Println("没有可用的AI客户端工具配置 MCP")
		return node.Router(requestParameter, dynamicContext)
	}

	// 遍历处理每个MCP配置
	for _, mcpVO := range aiClientToolMcpList {
		// 创建McpSyncClient对象
		mcpSyncClient, err := node.createMcpSyncClient(mcpVO)
		if err != nil {
			log.Printf("创建MCP客户端失败: %v", err)
			continue
		}

		// 注册Bean
		beanName := node.beanName(mcpVO.ID)
		node.RegisterDependency(beanName, mcpSyncClient)
	}

	return node.Router(requestParameter, dynamicContext)
}

// Get 获取下一个处理器
func (node *AiClientToolMcpNode) Get(requestParameter *entity.AiAgentEngineStarterEntity, dynamicContext *context.DynamicContext) (StrategyHandler, error) {
	return node.AiClientAdvisorNode, nil
}

// Router 路由到下一个处理器
func (node *AiClientToolMcpNode) Router(requestParameter *entity.AiAgentEngineStarterEntity, dynamicContext *context.DynamicContext) (string, error) {
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
func (node *AiClientToolMcpNode) beanName(id int64) string {
	return "AiClientToolMcp_" + strconv.FormatInt(id, 10)
}

// createMcpSyncClient 创建MCP同步客户端
func (node *AiClientToolMcpNode) createMcpSyncClient(aiClientToolMcpVO valobj.AiClientToolMcpVO) (McpSyncClient, error) {
	transportType := aiClientToolMcpVO.TransportType

	switch transportType {
	case "sse":
		return node.createSseMcpClient(aiClientToolMcpVO)
	case "stdio":
		return node.createStdioMcpClient(aiClientToolMcpVO)
	default:
		return nil, fmt.Errorf("err! transportType %s not exist!", transportType)
	}
}

// createSseMcpClient 创建SSE MCP客户端
func (node *AiClientToolMcpNode) createSseMcpClient(aiClientToolMcpVO valobj.AiClientToolMcpVO) (McpSyncClient, error) {
	transportConfigSse := aiClientToolMcpVO.TransportConfigSse
	if transportConfigSse == nil {
		return nil, errors.New("SSE传输配置为空")
	}

	originalBaseURI := transportConfigSse.BaseURI
	var baseURI, sseEndpoint string

	// 解析URL，处理查询参数
	if strings.Contains(originalBaseURI, "sse") {
		queryParamStartIndex := strings.Index(originalBaseURI, "sse")
		if queryParamStartIndex > 0 {
			baseURI = originalBaseURI[:queryParamStartIndex-1]
			sseEndpoint = originalBaseURI[queryParamStartIndex-1:]
		} else {
			baseURI = originalBaseURI
			sseEndpoint = transportConfigSse.SseEndpoint
		}
	} else {
		baseURI = originalBaseURI
		sseEndpoint = transportConfigSse.SseEndpoint
	}

	if sseEndpoint == "" {
		sseEndpoint = "/sse"
	}

	// 创建SSE传输客户端
	sseClientTransport := &HttpClientSseClientTransport{
		BaseURI:     baseURI,
		SseEndpoint: sseEndpoint,
	}

	// 创建MCP客户端
	mcpSyncClient := &MockMcpSyncClient{
		transport:      sseClientTransport,
		requestTimeout: time.Duration(aiClientToolMcpVO.RequestTimeout) * time.Minute,
	}

	// 初始化客户端
	initResult, err := mcpSyncClient.Initialize()
	if err != nil {
		return nil, fmt.Errorf("SSE MCP初始化失败: %v", err)
	}

	log.Printf("Tool SSE MCP Initialized %+v", initResult)
	return mcpSyncClient, nil
}

// createStdioMcpClient 创建Stdio MCP客户端
func (node *AiClientToolMcpNode) createStdioMcpClient(aiClientToolMcpVO valobj.AiClientToolMcpVO) (McpSyncClient, error) {
	transportConfigStdio := aiClientToolMcpVO.TransportConfigStdio
	if transportConfigStdio == nil {
		return nil, errors.New("Stdio传输配置为空")
	}

	stdioMap := transportConfigStdio.Stdio
	stdio, exists := stdioMap[aiClientToolMcpVO.McpName]
	if !exists {
		return nil, fmt.Errorf("找不到MCP名称 %s 对应的Stdio配置", aiClientToolMcpVO.McpName)
	}

	// 创建服务器参数
	stdioParams := &ServerParameters{
		Command: stdio.Command,
		Args:    stdio.Args,
	}

	// 创建Stdio传输客户端
	stdioClientTransport := &StdioClientTransport{
		ServerParams: stdioParams,
	}

	// 创建MCP客户端
	mcpSyncClient := &MockMcpSyncClient{
		transport:      stdioClientTransport,
		requestTimeout: time.Duration(aiClientToolMcpVO.RequestTimeout) * time.Second,
	}

	// 初始化客户端
	initResult, err := mcpSyncClient.Initialize()
	if err != nil {
		return nil, fmt.Errorf("Stdio MCP初始化失败: %v", err)
	}

	log.Printf("Tool Stdio MCP Initialized %+v", initResult)
	return mcpSyncClient, nil
}
