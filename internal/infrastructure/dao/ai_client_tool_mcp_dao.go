package dao

import (
	"time"

	"gorm.io/gorm"
	"smart-weaver/internal/infrastructure/dao/po"
)

// AiClientToolMcpDao MCP客户端配置数据访问对象
type AiClientToolMcpDao struct {
	DB *gorm.DB
}

// QueryAllMcpConfig 查询所有MCP配置
func (dao *AiClientToolMcpDao) QueryAllMcpConfig() ([]po.AiClientToolMcp, error) {
	var result []po.AiClientToolMcp
	if err := dao.DB.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

// QueryMcpConfigById 根据ID查询MCP配置
func (dao *AiClientToolMcpDao) QueryMcpConfigById(id int64) (*po.AiClientToolMcp, error) {
	var m po.AiClientToolMcp
	if err := dao.DB.First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// QueryMcpConfigByName 根据MCP名称查询配置
func (dao *AiClientToolMcpDao) QueryMcpConfigByName(name string) (*po.AiClientToolMcp, error) {
	var m po.AiClientToolMcp
	if err := dao.DB.Where("mcp_name = ?", name).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// Insert 插入MCP配置
func (dao *AiClientToolMcpDao) Insert(m *po.AiClientToolMcp) error {
	now := time.Now()
	m.CreateTime = now
	m.UpdateTime = now
	return dao.DB.Create(m).Error
}

// Update 更新MCP配置
func (dao *AiClientToolMcpDao) Update(m *po.AiClientToolMcp) error {
	m.UpdateTime = time.Now()
	return dao.DB.Save(m).Error
}

// DeleteById 根据ID删除MCP配置
func (dao *AiClientToolMcpDao) DeleteById(id int64) error {
	return dao.DB.Delete(&po.AiClientToolMcp{}, id).Error
}

// QueryMcpConfigByClientIds 根据客户端ID列表查询MCP配置
func (dao *AiClientToolMcpDao) QueryMcpConfigByClientIds(clientIds []int64) ([]po.AiClientToolMcp, error) {
	if len(clientIds) == 0 {
		return nil, nil
	}

	var result []po.AiClientToolMcp
	if err := dao.DB.Where("client_id IN ?", clientIds).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
