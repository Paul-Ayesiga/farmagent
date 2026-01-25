package repository

import (
	"github.com/farmagent/fa-crop-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HealthRecordRepository interface {
	Create(record *models.HealthRecord) error
	FindByID(id uuid.UUID) (*models.HealthRecord, error)
	FindByCropID(cropID uuid.UUID) ([]models.HealthRecord, error)
	FindLatestByCropID(cropID uuid.UUID) (*models.HealthRecord, error)
	Update(record *models.HealthRecord) error
	Delete(id uuid.UUID) error
}

type healthRecordRepository struct {
	db *gorm.DB
}

func NewHealthRecordRepository(db *gorm.DB) HealthRecordRepository {
	return &healthRecordRepository{db: db}
}

func (r *healthRecordRepository) Create(record *models.HealthRecord) error {
	return r.db.Create(record).Error
}

func (r *healthRecordRepository) FindByID(id uuid.UUID) (*models.HealthRecord, error) {
	var record models.HealthRecord
	err := r.db.Preload("Crop").Preload("Treatments").First(&record, "id = ?", id).Error
	return &record, err
}

func (r *healthRecordRepository) FindByCropID(cropID uuid.UUID) ([]models.HealthRecord, error) {
	var records []models.HealthRecord
	err := r.db.Where("crop_id = ?", cropID).Order("check_date DESC").Find(&records).Error
	return records, err
}

func (r *healthRecordRepository) FindLatestByCropID(cropID uuid.UUID) (*models.HealthRecord, error) {
	var record models.HealthRecord
	err := r.db.Where("crop_id = ?", cropID).Order("check_date DESC").First(&record).Error
	return &record, err
}

func (r *healthRecordRepository) Update(record *models.HealthRecord) error {
	return r.db.Save(record).Error
}

func (r *healthRecordRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.HealthRecord{}, "id = ?", id).Error
}
