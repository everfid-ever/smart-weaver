package dao

import (
	"gorm.io/gorm"
	"smart-weaver/internal/infrastructure/dao/po"
)

// IAiClientModelDao AI模型配置数据访问接口
type IAiClientModelDao interface {
	QueryAllModelConfig() ([]*po.AiClientModel, error)
	QueryModelConfigById(id uint) (*po.AiClientModel, error)
	QueryModelConfigByName(modelName string) (*po.AiClientModel, error)
	Insert(model *po.AiClientModel) error
	Update(model *po.AiClientModel) error
	DeleteById(id uint) error
}

type aiClientModelDao struct {
	db *gorm.DB
}

// NewAiClientModelDao 创建AI模型配置DAO
func NewAiClientModelDao(db *gorm.DB) IAiClientModelDao {
	return &aiClientModelDao{db: db}
}

// QueryAllModelConfig 查询所有模型配置
func (d *aiClientModelDao) QueryAllModelConfig() ([]*po.AiClientModel, error) {
	var models []*po.AiClientModel
	err := d.db.Order("id").Find(&models).Error
	return models, err
}

// QueryModelConfigById 根据ID查询模型配置
func (d *aiClientModelDao) QueryModelConfigById(id uint) (*po.AiClientModel, error) {
	var model po.AiClientModel
	err := d.db.First(&model, id).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

// QueryModelConfigByName 根据模型名称查询模型配置
func (d *aiClientModelDao) QueryModelConfigByName(modelName string) (*po.AiClientModel, error) {
	var model po.AiClientModel
	err := d.db.Where("model_name = ?", modelName).First(&model).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

// Insert 插入模型配置
func (d *aiClientModelDao) Insert(model *po.AiClientModel) error {
	return d.db.Create(model).Error
}

// Update 更新模型配置
func (d *aiClientModelDao) Update(model *po.AiClientModel) error {
	return d.db.Save(model).Error
}

// DeleteById 根据ID删除模型配置
func (d *aiClientModelDao) DeleteById(id uint) error {
	return d.db.Delete(&po.AiClientModel{}, id).Error
}
