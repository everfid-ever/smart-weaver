package po

import "time"

// AiClientModel defines the AI client model configuration
type AiClientModel struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ModelName       string    `gorm:"size:50;not null;uniqueIndex" json:"model_name"`
	BaseUrl         string    `gorm:"size:255;not null" json:"base_url"`
	ApiKey          string    `gorm:"size:255;not null" json:"api_key"`
	CompletionsPath string    `gorm:"size:100;default:v1/chat/completions" json:"completions_path"`
	EmbeddingsPath  string    `gorm:"size:100;default:v1/embeddings" json:"embeddings_path"`
	ModelType       string    `gorm:"size:50;not null" json:"model_type"`
	ModelVersion    string    `gorm:"size:50;default:gpt-4.1-mini" json:"model_version"`
	Timeout         int       `gorm:"default:30" json:"timeout"`
	Status          int       `gorm:"default:1" json:"status"`
	CreateTime      time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime      time.Time `gorm:"autoUpdateTime" json:"update_time"`
}

// TableName returns table name
func (AiClientModel) TableName() string {
	return "ai_client_model"
}
