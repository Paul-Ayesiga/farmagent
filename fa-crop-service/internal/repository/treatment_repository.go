package repository

import (
	"github.com/farmagent/fa-crop-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TreatmentRepository interface {
	Create(treatment *models.Treatment) error
	FindByID(id uuid.UUID) (*models.Treatment, error)
	FindByCropID(cropID uuid.UUID) ([]models.Treatment, error)
	FindByHealthRecordID(healthRecordID uuid.UUID) ([]models.Treatment, error)
	Update(treatment *models.Treatment) error
	Delete(id uuid.UUID) error
}

type treatmentRepository struct {
	db *gorm.DB
}

func NewTreatmentRepository(db *gorm.DB) TreatmentRepository {
	return &treatmentRepository{db: db}
}

func (r *treatmentRepository) Create(treatment *models.Treatment) error {
	return r.db.Create(treatment).Error
}

func (r *treatmentRepository) FindByID(id uuid.UUID) (*models.Treatment, error) {
	var treatment models.Treatment
	err := r.db.Preload("Crop").First(&treatment, "id = ?", id).Error
	return &treatment, err
}

func (r *treatmentRepository) FindByCropID(cropID uuid.UUID) ([]models.Treatment, error) {
	var treatments []models.Treatment
	err := r.db.Where("crop_id = ?", cropID).Order("application_date DESC").Find(&treatments).Error
	return treatments, err
}

func (r *treatmentRepository) FindByHealthRecordID(healthRecordID uuid.UUID) ([]models.Treatment, error) {
	var treatments []models.Treatment
	err := r.db.Where("health_record_id = ?", healthRecordID).Find(&treatments).Error
	return treatments, err
}

func (r *treatmentRepository) Update(treatment *models.Treatment) error {
	return r.db.Save(treatment).Error
}

func (r *treatmentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Treatment{}, "id = ?", id).Error
}
