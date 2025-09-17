package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"smart-weaver/internal/infrastructure/dao/po"
)

// InitDatabase 初始化数据库
func InitDatabase(cfg *Config) *gorm.DB {
	// MySQL DSN格式: username:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Datasource.Username,
		cfg.Database.Datasource.Password,
		cfg.Database.Datasource.Database,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(&po.AiClientModel{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}
