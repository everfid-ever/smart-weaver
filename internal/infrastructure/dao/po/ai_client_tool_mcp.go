package po

import "time"

// AiClientToolMcp MCP客户端配置表
type AiClientToolMcp struct {
	// 主键ID
	ID int64 `json:"id"`

	// MCP名称
	McpName string `json:"mcp_name"`

	// 传输类型(sse/stdio)
	TransportType string `json:"transport_type"`

	// 传输配置
	TransportConfig string `json:"transport_config"`

	// 请求超时时间(分钟)
	RequestTimeout int `json:"request_timeout"`

	// 状态(0:禁用,1:启用)
	Status int `json:"status"`

	// 创建时间
	CreateTime time.Time `json:"create_time"`

	// 更新时间
	UpdateTime time.Time `json:"update_time"`
}
