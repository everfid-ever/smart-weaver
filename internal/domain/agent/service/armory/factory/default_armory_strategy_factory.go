package factory

import (
	"smart-weaver/internal/domain/agent/service/armory/node"
)

// DefaultArmoryStrategyFactory 工厂类
type DefaultArmoryStrategyFactory struct {
	rootNode *node.RootNode
}

// NewDefaultArmoryStrategyFactory 创建工厂实例
func NewDefaultArmoryStrategyFactory(rootNode *node.RootNode) *DefaultArmoryStrategyFactory {
	return &DefaultArmoryStrategyFactory{
		rootNode: rootNode,
	}
}

// StrategyHandler 返回策略处理器
func (f *DefaultArmoryStrategyFactory) StrategyHandler() node.StrategyHandler {
	return f.rootNode
}
