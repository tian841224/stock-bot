package repository

import (
	"stock-bot/internal/db/models"

	"gorm.io/gorm"
)

type FeatureRepository interface {
	Create(feature *models.Feature) error
	GetByID(id uint) (*models.Feature, error)
	GetByCode(code string) (*models.Feature, error)
	GetByName(name string) (*models.Feature, error)
	Update(feature *models.Feature) error
	Delete(id uint) error
	List(offset, limit int) ([]*models.Feature, error)
	GetAll() ([]*models.Feature, error)
}

type featureRepository struct {
	db *gorm.DB
}

func NewFeatureRepository(db *gorm.DB) FeatureRepository {
	return &featureRepository{db: db}
}

// Create 建立新功能
func (r *featureRepository) Create(feature *models.Feature) error {
	return r.db.Create(feature).Error
}

// GetByID 根據 ID 取得功能
func (r *featureRepository) GetByID(id uint) (*models.Feature, error) {
	var feature models.Feature
	err := r.db.First(&feature, id).Error
	if err != nil {
		return nil, err
	}
	return &feature, nil
}

// GetByCode 根據代碼取得功能
func (r *featureRepository) GetByCode(code string) (*models.Feature, error) {
	var feature models.Feature
	err := r.db.Where("code = ?", code).First(&feature).Error
	if err != nil {
		return nil, err
	}
	return &feature, nil
}

// GetByName 根據名稱取得功能
func (r *featureRepository) GetByName(name string) (*models.Feature, error) {
	var feature models.Feature
	err := r.db.Where("name = ?", name).First(&feature).Error
	if err != nil {
		return nil, err
	}
	return &feature, nil
}

// Update 更新功能資料
func (r *featureRepository) Update(feature *models.Feature) error {
	return r.db.Save(feature).Error
}

// Delete 刪除功能
func (r *featureRepository) Delete(id uint) error {
	return r.db.Delete(&models.Feature{}, id).Error
}

// List 取得功能列表
func (r *featureRepository) List(offset, limit int) ([]*models.Feature, error) {
	var features []*models.Feature
	err := r.db.Offset(offset).Limit(limit).Find(&features).Error
	return features, err
}

// GetAll 取得所有功能
func (r *featureRepository) GetAll() ([]*models.Feature, error) {
	var features []*models.Feature
	err := r.db.Find(&features).Error
	return features, err
}
