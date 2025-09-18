package context

import "sync"

// DynamicContext 动态上下文
type DynamicContext struct {
	Level       int                    `json:"level"`
	DataObjects map[string]interface{} `json:"data_objects"`
	mu          sync.RWMutex           // 读写锁保护并发访问
}

// NewDynamicContext 创建动态上下文实例
func NewDynamicContext() *DynamicContext {
	return &DynamicContext{
		Level:       0,
		DataObjects: make(map[string]interface{}),
	}
}

// NewDynamicContextWithLevel 创建带级别的动态上下文实例
func NewDynamicContextWithLevel(level int) *DynamicContext {
	return &DynamicContext{
		Level:       level,
		DataObjects: make(map[string]interface{}),
	}
}

// SetValue 设置值
func (dc *DynamicContext) SetValue(key string, value interface{}) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	dc.DataObjects[key] = value
}

// GetValue 获取值
func (dc *DynamicContext) GetValue(key string) interface{} {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return dc.DataObjects[key]
}

// GetValueWithType 获取指定类型的值
func (dc *DynamicContext) GetValueWithType(key string) (interface{}, bool) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	value, exists := dc.DataObjects[key]
	return value, exists
}

// GetLevel 获取级别
func (dc *DynamicContext) GetLevel() int {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return dc.Level
}

// SetLevel 设置级别
func (dc *DynamicContext) SetLevel(level int) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	dc.Level = level
}

// HasKey 检查是否存在指定key
func (dc *DynamicContext) HasKey(key string) bool {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	_, exists := dc.DataObjects[key]
	return exists
}

// RemoveKey 移除指定key
func (dc *DynamicContext) RemoveKey(key string) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	delete(dc.DataObjects, key)
}

// Clear 清空所有数据
func (dc *DynamicContext) Clear() {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	dc.DataObjects = make(map[string]interface{})
}

// Keys 获取所有key
func (dc *DynamicContext) Keys() []string {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	keys := make([]string, 0, len(dc.DataObjects))
	for key := range dc.DataObjects {
		keys = append(keys, key)
	}
	return keys
}

// Size 获取数据对象数量
func (dc *DynamicContext) Size() int {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return len(dc.DataObjects)
}
