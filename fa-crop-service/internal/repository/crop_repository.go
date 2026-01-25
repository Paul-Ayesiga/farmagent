package repository

import (
	"github.com/farmagent/fa-crop-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CropRepository interface {
	Create(crop *models.Crop) error
	FindByID(id uuid.UUID) (*models.Crop, error)
	FindByFieldID(fieldID uuid.UUID) ([]models.Crop, error)
	FindByUserID(userID string) ([]models.Crop, error)
	Update(crop *models.Crop) error
	Delete(id uuid.UUID) error
}

type cropRepository struct {
	db *gorm.DB
}

func NewCropRepository(db *gorm.DB) CropRepository {
	return &cropRepository{db: db}
}

func (r *cropRepository) Create(crop *models.Crop) error {
	return r.db.Create(crop).Error
}

func (r *cropRepository) FindByID(id uuid.UUID) (*models.Crop, error) {
	var crop models.Crop
	err := r.db.Preload("Field").Preload("HealthRecords").Preload("Treatments").
		First(&crop, "id = ?", id).Error
	return &crop, err
}

func (r *cropRepository) FindByFieldID(fieldID uuid.UUID) ([]models.Crop, error) {
	var crops []models.Crop
	err := r.db.Where("field_id = ?", fieldID).Order("planting_date DESC").Find(&crops).Error
	return crops, err
}

func (r *cropRepository) FindByUserID(userID string) ([]models.Crop, error) {
	var crops []models.Crop
	err := r.db.Joins("JOIN fields ON fields.id = crops.field_id").
		Where("fields.user_id = ?", userID).
		Preload("Field").
		Order("crops.created_at DESC").
		Find(&crops).Error
	return crops, err
}

func (r *cropRepository) Update(crop *models.Crop) error {
	return r.db.Save(crop).Error
}

func (r *cropRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Crop{}, "id = ?", id).Error
}
