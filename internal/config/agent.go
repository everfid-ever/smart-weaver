package config

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL 驱动
	_ "github.com/lib/pq"              // PostgreSQL 驱动
)

// AiAgentConfig AI代理配置
type AiAgentConfig struct {
	// 主数据库配置
	MainDB struct {
		Driver       string `yaml:"driver"`
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		Database     string `yaml:"database"`
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
		MaxOpenConns int    `yaml:"max_open_conns" default:"10"`
		MaxIdleConns int    `yaml:"max_idle_conns" default:"5"`
		MaxLifetime  int    `yaml:"max_lifetime" default:"1800"` // 秒
		MaxIdleTime  int    `yaml:"max_idle_time" default:"30"`  // 秒
	} `yaml:"main_db"`

	// PgVector 数据库配置
	VectorDB struct {
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		Database     string `yaml:"database"`
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
		MaxOpenConns int    `yaml:"max_open_conns" default:"5"`
		MaxIdleConns int    `yaml:"max_idle_conns" default:"2"`
		MaxIdleTime  int    `yaml:"max_idle_time" default:"30"` // 秒
	} `yaml:"vector_db"`

	// OpenAI 配置
	OpenAI struct {
		BaseUrl string `yaml:"base_url"`
		ApiKey  string `yaml:"api_key"`
	} `yaml:"openai"`
}

// DataSource 数据源
type DataSource struct {
	*sql.DB
	Name string
}

// OpenAiApi OpenAI API 客户端
type OpenAiApi struct {
	BaseUrl string
	ApiKey  string
}

// OpenAiEmbeddingModel OpenAI 嵌入模型
type OpenAiEmbeddingModel struct {
	Api *OpenAiApi
}

// PgVectorStore PG向量存储
type PgVectorStore struct {
	DB              *sql.DB
	EmbeddingModel  *OpenAiEmbeddingModel
	VectorTableName string
}

// TokenTextSplitter 文本分割器
type TokenTextSplitter struct{}

// MainDataSource 创建主数据源（替代MyBatis数据源）
func (config *AiAgentConfig) MainDataSource() (*DataSource, error) {
	var dsn string

	switch config.MainDB.Driver {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4",
			config.MainDB.Username,
			config.MainDB.Password,
			config.MainDB.Host,
			config.MainDB.Port,
			config.MainDB.Database)
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.MainDB.Host,
			config.MainDB.Port,
			config.MainDB.Username,
			config.MainDB.Password,
			config.MainDB.Database)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", config.MainDB.Driver)
	}

	db, err := sql.Open(config.MainDB.Driver, dsn)
	if err != nil {
		return nil, err
	}

	// 连接池配置
	db.SetMaxOpenConns(config.MainDB.MaxOpenConns)
	db.SetMaxIdleConns(config.MainDB.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(config.MainDB.MaxLifetime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(config.MainDB.MaxIdleTime) * time.Second)

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DataSource{
		DB:   db,
		Name: "MainDataSource",
	}, nil
}

// PgVectorDataSource 创建 PgVector 数据源
func (config *AiAgentConfig) PgVectorDataSource() (*DataSource, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.VectorDB.Host,
		config.VectorDB.Port,
		config.VectorDB.Username,
		config.VectorDB.Password,
		config.VectorDB.Database)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// 连接池配置
	db.SetMaxOpenConns(config.VectorDB.MaxOpenConns)
	db.SetMaxIdleConns(config.VectorDB.MaxIdleConns)
	db.SetConnMaxIdleTime(time.Duration(config.VectorDB.MaxIdleTime) * time.Second)

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to PgVector database: %w", err)
	}

	return &DataSource{
		DB:   db,
		Name: "PgVectorDataSource",
	}, nil
}

// VectorStore 创建向量存储
func (config *AiAgentConfig) VectorStore() (*PgVectorStore, error) {
	// 创建 PgVector 数据源
	pgVectorDS, err := config.PgVectorDataSource()
	if err != nil {
		return nil, err
	}

	// 创建 OpenAI API 客户端
	openAiApi := &OpenAiApi{
		BaseUrl: config.OpenAI.BaseUrl,
		ApiKey:  config.OpenAI.ApiKey,
	}

	// 创建嵌入模型
	embeddingModel := &OpenAiEmbeddingModel{
		Api: openAiApi,
	}

	return &PgVectorStore{
		DB:              pgVectorDS.DB,
		EmbeddingModel:  embeddingModel,
		VectorTableName: "vector_store_openai",
	}, nil
}

// TokenTextSplitter 创建文本分割器
func (config *AiAgentConfig) CreateTokenTextSplitter() *TokenTextSplitter {
	return &TokenTextSplitter{}
}

// Close 关闭数据源
func (ds *DataSource) Close() error {
	if ds.DB != nil {
		return ds.DB.Close()
	}
	return nil
}
