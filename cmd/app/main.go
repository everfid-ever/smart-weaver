package main

import (
	"log"

	"smart-weaver/internal/config"
	"smart-weaver/internal/trigger/http"
)

func main() {
	// 初始化配置
	cfg := config.Load()

	// 初始化数据库
	db := config.InitDatabase(cfg)

	// 初始化缓存
	cache := config.InitCache()

	// 初始化线程池
	threadPool := config.InitThreadPool(cfg)

	// 启动HTTP服务器
	router := http.SetupRouter(db, cache, threadPool)

	port := cfg.Server.Port
	if port == "" {
		port = "8091"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
