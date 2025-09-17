package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"spring"`
	ThreadPool ThreadPoolConfig `mapstructure:"thread"`
	Logging    LoggingConfig    `mapstructure:"logging"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Datasource DatasourceConfig `mapstructure:"datasource"`
}

// DatasourceConfig 数据源配置
type DatasourceConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Driver   string `mapstructure:"driver-class-name"`
}

// ThreadPoolConfig 线程池配置
type ThreadPoolConfig struct {
	Pool PoolConfig `mapstructure:"pool"`
}

// PoolConfig 池配置
type PoolConfig struct {
	Executor ExecutorConfig `mapstructure:"executor"`
}

// ExecutorConfig 执行器配置
type ExecutorConfig struct {
	Config ExecutorConfigDetail `mapstructure:"config"`
}

// ExecutorConfigDetail 执行器配置详情
type ExecutorConfigDetail struct {
	CorePoolSize   int    `mapstructure:"core-pool-size"`
	MaxPoolSize    int    `mapstructure:"max-pool-size"`
	KeepAliveTime  int64  `mapstructure:"keep-alive-time"`
	BlockQueueSize int    `mapstructure:"block-queue-size"`
	Policy         string `mapstructure:"policy"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level  map[string]string `mapstructure:"level"`
	Config string            `mapstructure:"config"`
}

// Load 加载配置
func Load() *Config {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	// 设置默认值
	viper.SetDefault("server.port", "8091")

	// 获取环境变量
	profile := viper.GetString("spring.profiles.active")
	if profile == "" {
		profile = "dev"
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
	}

	// 合并环境特定配置
	viper.SetConfigName("application-" + profile)
	if err := viper.MergeInConfig(); err != nil {
		log.Printf("Error merging config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	return &config
}
