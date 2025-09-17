package config

import (
	"time"
	
	"github.com/patrickmn/go-cache"
)

// InitCache 初始化缓存
func InitCache() *cache.Cache {
	return cache.New(3*time.Second, 5*time.Minute)
}
