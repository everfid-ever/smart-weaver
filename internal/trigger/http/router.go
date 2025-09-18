package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"smart-weaver/internal/config"
)

// SetupRouter 设置路由
func SetupRouter(db *gorm.DB, cache *cache.Cache, threadPool *config.ThreadPoolExecutor) *gin.Engine {
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
		// 测试接口
		api.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "AI Agent Station Go服务运行正常",
			})
		})

		// 异步任务提交接口
		api.POST("/task", func(c *gin.Context) {
			taskData := c.PostForm("data")
			if taskData == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "data不能为空"})
				return
			}

			// 提交异步任务
			threadPool.Submit(func() {
				// 写入缓存，使用5分钟过期
				cache.Set(taskData, "processed", 5*time.Minute)
			})

			c.JSON(http.StatusOK, gin.H{"status": "submitted"})
		})

		// 缓存查询接口
		api.GET("/cache/:key", func(c *gin.Context) {
			key := c.Param("key")
			value, found := cache.Get(key)
			if !found {
				c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
		})
	}

	return router
}
