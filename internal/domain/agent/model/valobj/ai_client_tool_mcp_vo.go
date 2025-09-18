package valobj

// AiClientToolMcpVO MCP VO 对象
type AiClientToolMcpVO struct {
	ID                   int64                 `json:"id"`
	McpName              string                `json:"mcp_name"`
	TransportType        string                `json:"transport_type"`         // sse / stdio
	TransportConfigSse   *TransportConfigSse   `json:"transport_config_sse"`   // SSE 配置，可为空
	TransportConfigStdio *TransportConfigStdio `json:"transport_config_stdio"` // STDIO 配置，可为空
	RequestTimeout       int                   `json:"request_timeout"`        // 分钟
}

// TransportConfigSse SSE 配置
type TransportConfigSse struct {
	BaseURI     string `json:"base_uri"`
	SseEndpoint string `json:"sse_endpoint"`
}

// TransportConfigStdio STDIO 配置
type TransportConfigStdio struct {
	Stdio map[string]Stdio `json:"stdio"` // key 对应 mcp-server 名称
}

// Stdio STDIO 命令配置
type Stdio struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}
