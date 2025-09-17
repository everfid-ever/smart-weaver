package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"smart-weaver/internal/config"
)

// SetupRouter 设置路由
func SetupRouter(db *gorm.DB, cache *cache.Cache, threadPool *config.ThreadPool) *gin.Engine {
	router := gin.Default()

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	// API路由组
	api := router.Group("/api/v1")
	{
		// 这里可以添加具体的业务路由
		api.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "AI Agent Station Go服务运行正常",
			})
		})
	}

	return router
}
