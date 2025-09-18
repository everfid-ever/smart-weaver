package po

import (
	"time"
	
	"smart-weaver/internal/infrastructure/dao/po/base"
)

type AiClientModel struct {
	base.Page

	ID              int64     `json:"id"`
	ModelName       string    `json:"model_name"`
	BaseURL         string    `json:"base_url"`
	APIKey          string    `json:"api_key"`
	CompletionsPath string    `json:"completions_path"`
	EmbeddingsPath  string    `json:"embeddings_path"`
	ModelType       string    `json:"model_type"`
	ModelVersion    string    `json:"model_version"`
	Timeout         int       `json:"timeout"`
	Status          int       `json:"status"`
	CreateTime      time.Time `json:"create_time"`
	UpdateTime      time.Time `json:"update_time"`
}
