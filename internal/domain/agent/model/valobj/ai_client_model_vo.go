package valobj

import "time"

type AiClientModelVO struct {
	ID                       int64                       `json:"id"`
	ModelName                string                      `json:"model_name"`
	BaseURL                  string                      `json:"base_url"`
	APIKey                   string                      `json:"api_key"`
	CompletionsPath          string                      `json:"completions_path"`
	EmbeddingsPath           string                      `json:"embeddings_path"`
	ModelType                string                      `json:"model_type"` // openai/azure 等
	ModelVersion             string                      `json:"model_version"`
	Timeout                  int                         `json:"timeout"` // 秒
	AIClientModelToolConfigs []AIClientModelToolConfigVO `json:"ai_client_model_tool_configs"`
}

// AIClientModelToolConfigVO 嵌套工具配置
type AIClientModelToolConfigVO struct {
	ID         int       `json:"id"`
	ModelID    int64     `json:"model_id"`
	ToolType   string    `json:"tool_type"` // mcp / function call
	ToolID     int64     `json:"tool_id"`   // MCP ID / function call ID
	CreateTime time.Time `json:"create_time"`
}
