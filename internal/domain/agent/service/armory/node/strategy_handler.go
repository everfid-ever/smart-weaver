package node

import (
	"smart-weaver/internal/domain/agent/model/entity"
	"smart-weaver/internal/domain/agent/service/armory/factory/context"
)

// StrategyHandler 策略处理器统一接口
type StrategyHandler interface {
	// DoApply 执行应用逻辑
	DoApply(requestParameter *entity.AiAgentEngineStarterEntity, dynamicContext *context.DynamicContext) (string, error)
	// Get 获取下一个处理器
	Get(requestParameter *entity.AiAgentEngineStarterEntity, dynamicContext *context.DynamicContext) (StrategyHandler, error)
	// Router 路由到下一个处理器
	Router(requestParameter *entity.AiAgentEngineStarterEntity, dynamicContext *context.DynamicContext) (string, error)
}
