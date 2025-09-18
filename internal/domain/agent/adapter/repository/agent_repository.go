package repository

import "smart-weaver/internal/domain/agent/model/valobj"

type IAgentRepository interface {
	// QueryAiClientModelVOListByClientIDs 根据 clientId 列表查询 AiClientModelVO
	QueryAiClientModelVOListByClientIDs(clientIDList []int64) ([]valobj.AiClientModelVO, error)

	// QueryAiClientToolMcpVOListByClientIDs 根据 clientId 列表查询 AiClientToolMcpVO
	QueryAiClientToolMcpVOListByClientIDs(clientIDList []int64) ([]valobj.AiClientToolMcpVO, error)
}
