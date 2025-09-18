package dao

import (
	"gorm.io/gorm"
	"smart-weaver/internal/infrastructure/dao/po"
)

type AiClientModelDao struct {
	DB *gorm.DB
}

// QueryAllModelConfig 查询所有模型配置
func (d *AiClientModelDao) QueryAllModelConfig() ([]po.AiClientModel, error) {
	var models []po.AiClientModel
	err := d.DB.Find(&models).Error
	return models, err
}

// QueryModelConfigById 根据ID查询模型配置
func (d *AiClientModelDao) QueryModelConfigById(id int64) (*po.AiClientModel, error) {
	var model po.AiClientModel
	err := d.DB.First(&model, id).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

// QueryModelConfigByName 根据模型名称查询模型配置
func (d *AiClientModelDao) QueryModelConfigByName(name string) (*po.AiClientModel, error) {
	var model po.AiClientModel
	err := d.DB.Where("model_name = ?", name).First(&model).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

// Insert 插入模型配置
func (d *AiClientModelDao) Insert(model *po.AiClientModel) error {
	return d.DB.Create(model).Error
}

// Update 更新模型配置
func (d *AiClientModelDao) Update(model *po.AiClientModel) error {
	return d.DB.Save(model).Error
}

// DeleteById 根据ID删除模型配置
func (d *AiClientModelDao) DeleteById(id int64) error {
	return d.DB.Delete(&po.AiClientModel{}, id).Error
}

// QueryModelConfigByClientIds 根据客户端ID列表查询模型配置
func (d *AiClientModelDao) QueryModelConfigByClientIds(clientIds []int64) ([]po.AiClientModel, error) {
	var models []po.AiClientModel
	err := d.DB.Where("id IN ?", clientIds).Find(&models).Error
	return models, err
}

// QueryToolMcpConfigByClientIds 根据客户端ID列表查询工具MCP配置
func (d *AiClientModelDao) QueryToolMcpConfigByClientIds(clientIds []int64) ([]po.AiClientToolMcp, error) {
	var tools []po.AiClientToolMcp
	err := d.DB.Where("id IN ?", clientIds).Find(&tools).Error
	return tools, err
}

// QueryClientModelList 根据条件查询客户端模型列表
func (d *AiClientModelDao) QueryClientModelList(filter *po.AiClientModel) ([]po.AiClientModel, error) {
	var models []po.AiClientModel
	query := d.DB.Model(&po.AiClientModel{})
	if filter.ModelName != "" {
		query = query.Where("model_name = ?", filter.ModelName)
	}
	if filter.Status != 0 {
		query = query.Where("status = ?", filter.Status)
	}
	err := query.Find(&models).Error
	return models, err
}
