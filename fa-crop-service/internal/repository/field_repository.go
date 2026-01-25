package repository

import (
	"github.com/farmagent/fa-crop-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FieldRepository interface {
	Create(field *models.Field) error
	FindByID(id uuid.UUID) (*models.Field, error)
	FindByUserID(userID string) ([]models.Field, error)
	Update(field *models.Field) error
	Delete(id uuid.UUID) error
}

type fieldRepository struct {
	db *gorm.DB
}

func NewFieldRepository(db *gorm.DB) FieldRepository {
	return &fieldRepository{db: db}
}

func (r *fieldRepository) Create(field *models.Field) error {
	return r.db.Create(field).Error
}

func (r *fieldRepository) FindByID(id uuid.UUID) (*models.Field, error) {
	var field models.Field
	err := r.db.Preload("Crops").First(&field, "id = ?", id).Error
	return &field, err
}

func (r *fieldRepository) FindByUserID(userID string) ([]models.Field, error) {
	var fields []models.Field
	err := r.db.Preload("Crops").Where("user_id = ?", userID).Order("created_at DESC").Find(&fields).Error
	return fields, err
}

func (r *fieldRepository) Update(field *models.Field) error {
	return r.db.Save(field).Error
}

func (r *fieldRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Field{}, "id = ?", id).Error
}
