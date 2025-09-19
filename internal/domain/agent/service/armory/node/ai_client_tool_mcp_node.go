package node

import (
	"fmt"
	"log"
	"smart-weaver/internal/domain/agent/model/entity"
	"smart-weaver/internal/domain/agent/model/valobj"
	"smart-weaver/internal/domain/agent/service/armory"
	"smart-weaver/internal/domain/agent/service/armory/factory/context"
)

// McpSyncClient MCP同步客户端结构体
type McpSyncClient struct {
	ID                   int64
	McpName              string
	TransportType        string                       // sse / stdio
	TransportConfigSse   *valobj.TransportConfigSse   // SSE 配置，可为空
	TransportConfigStdio *valobj.TransportConfigStdio // STDIO 配置，可为空
	RequestTimeout       int                          // 分钟
}

// McpClientRegistry MCP客户端注册器接口
type McpClientRegistry interface {
	RegisterClient(beanName string, client *McpSyncClient) error
	RemoveClient(beanName string) error
	ContainsClient(beanName string) bool
	GetClient(beanName string) *McpSyncClient
}

// DefaultMcpClientRegistry 默认MCP客户端注册器实现
type DefaultMcpClientRegistry struct {
	clients map[string]*McpSyncClient
}

// NewDefaultMcpClientRegistry 创建默认MCP客户端注册器
func NewDefaultMcpClientRegistry() *DefaultMcpClientRegistry {
	return &DefaultMcpClientRegistry{
		clients: make(map[string]*McpSyncClient),
	}
}

// RegisterClient 注册客户端
func (r *DefaultMcpClientRegistry) RegisterClient(beanName string, client *McpSyncClient) error {
	r.clients[beanName] = client
	return nil
}

// RemoveClient 移除客户端
func (r *DefaultMcpClientRegistry) RemoveClient(beanName string) error {
	delete(r.clients, beanName)
	return nil
}

// ContainsClient 检查客户端是否存在
func (r *DefaultMcpClientRegistry) ContainsClient(beanName string) bool {
	_, exists := r.clients[beanName]
	return exists
}

// GetClient 获取客户端
func (r *DefaultMcpClientRegistry) GetClient(beanName string) *McpSyncClient {
	return r.clients[beanName]
}

// AiClientToolMcpNode AI客户端工具MCP节点
type AiClientToolMcpNode struct {
	*armory.AbstractArmorySupport
	mcpRegistry McpClientRegistry
}

// NewAiClientToolMcpNode 创建AI客户端工具MCP节点
func NewAiClientToolMcpNode() *AiClientToolMcpNode {
	return &AiClientToolMcpNode{
		AbstractArmorySupport: &armory.AbstractArmorySupport{
			ThreadPool: make(chan func(), 100), // 创建容量为100的线程池
			Deps:       make(map[string]any),
		},
		mcpRegistry: NewDefaultMcpClientRegistry(),
	}
}

// DoApply 执行应用逻辑
func (n *AiClientToolMcpNode) DoApply(requestParameter *entity.AiAgentEngineStarterEntity, dynamicContext *context.DynamicContext) (string, error) {
	log.Println("AiAgent 装配，tool mcp")

	// 从动态上下文获取AI客户端工具MCP列表
	aiClientToolMcpListInterface := dynamicContext.GetValue("aiClientToolMcpList")
	if aiClientToolMcpListInterface == nil {
		log.Println("warn: 没有可用的AI客户端工具配置 MCP")
		// 没有MCP配置也不是错误，继续执行
		return n.Router(requestParameter, dynamicContext)
	}

	aiClientToolMcpList, ok := aiClientToolMcpListInterface.([]valobj.AiClientToolMcpVO)
	if !ok || len(aiClientToolMcpList) == 0 {
		log.Println("warn: 没有可用的AI客户端工具配置 MCP")
		// 没有MCP配置也不是错误，继续执行
		return n.Router(requestParameter, dynamicContext)
	}

	// 遍历MCP列表，为每个MCP创建对应的客户端
	for _, mcpVO := range aiClientToolMcpList {
		// 构建Bean名称
		beanName := fmt.Sprintf("AiClientToolMcp%d", mcpVO.ID)

		// 创建McpSyncClient对象
		mcpSyncClient, err := n.createMcpSyncClient(mcpVO)
		if err != nil {
			log.Printf("error: 创建McpSyncClient失败: %v", err)
			continue
		}

		// 如果客户端已存在，先移除
		if n.mcpRegistry.ContainsClient(beanName) {
			if err := n.mcpRegistry.RemoveClient(beanName); err != nil {
				log.Printf("error: 移除现有MCP客户端失败: %v", err)
			}
		}

		// 注册新的客户端
		if err := n.mcpRegistry.RegisterClient(beanName, mcpSyncClient); err != nil {
			log.Printf("error: 注册MCP客户端失败: %v", err)
			continue
		}

		log.Printf("成功注册AI客户端工具MCP Bean: %s", beanName)
	}

	return n.Router(requestParameter, dynamicContext)
}

// Get 获取下一个处理器（原Java代码返回null，表示这是处理链的最后一个节点）
func (n *AiClientToolMcpNode) Get(requestParameter *entity.AiAgentEngineStarterEntity, dynamicContext *context.DynamicContext) (StrategyHandler, error) {
	return nil, nil
}

// Router 路由方法
func (n *AiClientToolMcpNode) Router(requestParameter *entity.AiAgentEngineStarterEntity, dynamicContext *context.DynamicContext) (string, error) {
	nextHandler, err := n.Get(requestParameter, dynamicContext)
	if err != nil {
		return "", err
	}
	if nextHandler != nil {
		log.Println("AiClientToolMcpNode 路由到下一个处理器")
		return nextHandler.DoApply(requestParameter, dynamicContext)
	}
	// 如果没有下一个处理器，返回成功标识
	log.Println("AiClientToolMcpNode 处理完成 - 这是处理链的最后一个节点")
	return "success", nil
}

// createMcpSyncClient 创建MCP同步客户端
func (n *AiClientToolMcpNode) createMcpSyncClient(mcpVO valobj.AiClientToolMcpVO) (*McpSyncClient, error) {
	if mcpVO.McpName == "" {
		return nil, fmt.Errorf("MCP name cannot be empty")
	}

	if mcpVO.TransportType == "" {
		return nil, fmt.Errorf("transport type cannot be empty")
	}

	// 验证传输类型
	if mcpVO.TransportType != "sse" && mcpVO.TransportType != "stdio" {
		return nil, fmt.Errorf("unsupported transport type: %s, must be 'sse' or 'stdio'", mcpVO.TransportType)
	}

	// 根据传输类型验证对应的配置
	switch mcpVO.TransportType {
	case "sse":
		if mcpVO.TransportConfigSse == nil {
			return nil, fmt.Errorf("SSE transport config cannot be nil when transport type is 'sse'")
		}
		if mcpVO.TransportConfigSse.BaseURI == "" {
			return nil, fmt.Errorf("SSE base URI cannot be empty")
		}
		if mcpVO.TransportConfigSse.SseEndpoint == "" {
			return nil, fmt.Errorf("SSE endpoint cannot be empty")
		}
	case "stdio":
		if mcpVO.TransportConfigStdio == nil {
			return nil, fmt.Errorf("STDIO transport config cannot be nil when transport type is 'stdio'")
		}
		if mcpVO.TransportConfigStdio.Stdio == nil || len(mcpVO.TransportConfigStdio.Stdio) == 0 {
			return nil, fmt.Errorf("STDIO config cannot be empty")
		}
		// 验证每个STDIO配置
		for serverName, stdio := range mcpVO.TransportConfigStdio.Stdio {
			if stdio.Command == "" {
				return nil, fmt.Errorf("STDIO command cannot be empty for server: %s", serverName)
			}
		}
	}

	mcpSyncClient := &McpSyncClient{
		ID:                   mcpVO.ID,
		McpName:              mcpVO.McpName,
		TransportType:        mcpVO.TransportType,
		TransportConfigSse:   mcpVO.TransportConfigSse,
		TransportConfigStdio: mcpVO.TransportConfigStdio,
		RequestTimeout:       mcpVO.RequestTimeout,
	}

	// 设置默认值
	if mcpSyncClient.RequestTimeout <= 0 {
		mcpSyncClient.RequestTimeout = 5 // 默认5分钟超时
	}

	return mcpSyncClient, nil
}

// GetMcpRegistry 获取MCP注册器（用于测试或外部访问）
func (n *AiClientToolMcpNode) GetMcpRegistry() McpClientRegistry {
	return n.mcpRegistry
}

// Connect 连接到MCP服务器（实际业务逻辑方法）
func (client *McpSyncClient) Connect() error {
	switch client.TransportType {
	case "sse":
		if client.TransportConfigSse == nil {
			return fmt.Errorf("SSE config is nil")
		}
		log.Printf("Connecting to MCP server via SSE: %s%s",
			client.TransportConfigSse.BaseURI, client.TransportConfigSse.SseEndpoint)
		// 这里应该实现实际的SSE连接逻辑
		return nil
	case "stdio":
		if client.TransportConfigStdio == nil {
			return fmt.Errorf("STDIO config is nil")
		}
		log.Printf("Connecting to MCP server via STDIO for %s", client.McpName)
		// 这里应该实现实际的STDIO连接逻辑
		for serverName, stdio := range client.TransportConfigStdio.Stdio {
			log.Printf("Starting STDIO server: %s with command: %s %v",
				serverName, stdio.Command, stdio.Args)
			// 实际的子进程启动逻辑应该在这里实现
		}
		return nil
	default:
		return fmt.Errorf("unsupported transport type: %s", client.TransportType)
	}
}

// Disconnect 断开MCP服务器连接
func (client *McpSyncClient) Disconnect() error {
	switch client.TransportType {
	case "sse":
		log.Printf("Disconnecting from MCP SSE server: %s", client.McpName)
		// 这里应该实现实际的SSE断连逻辑
		return nil
	case "stdio":
		log.Printf("Disconnecting from MCP STDIO servers for: %s", client.McpName)
		// 这里应该实现实际的STDIO进程终止逻辑
		for serverName := range client.TransportConfigStdio.Stdio {
			log.Printf("Stopping STDIO server: %s", serverName)
			// 实际的子进程终止逻辑应该在这里实现
		}
		return nil
	default:
		return fmt.Errorf("unsupported transport type: %s", client.TransportType)
	}
}

// IsConnected 检查是否已连接
func (client *McpSyncClient) IsConnected() bool {
	// 这里应该实现实际的连接状态检查
	// 目前只是一个占位符实现
	return true
}
