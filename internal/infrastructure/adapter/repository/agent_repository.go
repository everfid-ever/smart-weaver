package repository

import (
	"encoding/json"
	"log"

	"smart-weaver/internal/domain/agent/model/valobj"
	"smart-weaver/internal/infrastructure/dao"
)

type AgentRepository struct {
	clientModelDao *dao.AiClientModelDao
}

func NewAgentRepository(dao *dao.AiClientModelDao) *AgentRepository {
	return &AgentRepository{
		clientModelDao: dao,
	}
}

// QueryAiClientModelVOListByClientIds 查询 AI Client Model VO 列表
func (r *AgentRepository) QueryAiClientModelVOListByClientIds(clientIdList []int64) []*valobj.AiClientModelVO {
	aiClientModels, err := r.clientModelDao.QueryModelConfigByClientIds(clientIdList)
	if err != nil {
		log.Printf("查询模型配置失败: %v", err)
		return nil
	}

	voList := make([]*valobj.AiClientModelVO, 0, len(aiClientModels))

	for _, m := range aiClientModels {
		vo := &valobj.AiClientModelVO{
			ID:              m.ID,
			ModelName:       m.ModelName,
			BaseURL:         m.BaseURL,
			APIKey:          m.APIKey,
			CompletionsPath: m.CompletionsPath,
			EmbeddingsPath:  m.EmbeddingsPath,
			ModelType:       m.ModelType,
			ModelVersion:    m.ModelVersion,
			Timeout:         m.Timeout,
		}
		voList = append(voList, vo)
	}
	return voList
}

// QueryAiClientToolMcpVOListByClientIds 查询 AI Client Tool MCP VO 列表
func (r *AgentRepository) QueryAiClientToolMcpVOListByClientIds(clientIdList []int64) []*valobj.AiClientToolMcpVO {
	aiClientToolMcps, err := r.clientModelDao.QueryToolMcpConfigByClientIds(clientIdList)
	if err != nil {
		log.Printf("查询 MCP 配置失败: %v", err)
		return nil
	}
	voList := make([]*valobj.AiClientToolMcpVO, 0, len(aiClientToolMcps))

	for _, m := range aiClientToolMcps {
		vo := &valobj.AiClientToolMcpVO{
			ID:             m.ID,
			McpName:        m.McpName,
			TransportType:  m.TransportType,
			RequestTimeout: m.RequestTimeout,
		}

		if m.TransportConfig != "" {
			switch m.TransportType {
			case "sse":
				var sse valobj.TransportConfigSse
				if err := json.Unmarshal([]byte(m.TransportConfig), &sse); err != nil {
					log.Printf("解析 SSE 配置失败: %v", err)
				} else {
					vo.TransportConfigSse = &sse
				}
			case "stdio":
				var stdio valobj.TransportConfigStdio
				if err := json.Unmarshal([]byte(m.TransportConfig), &stdio); err != nil {
					log.Printf("解析 STDIO 配置失败: %v", err)
				} else {
					vo.TransportConfigStdio = &stdio
				}
			}
		}

		voList = append(voList, vo)
	}
	return voList
}
