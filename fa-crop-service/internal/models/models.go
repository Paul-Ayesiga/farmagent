package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Field represents a farmer's field/plot of land
type Field struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"type:varchar(100);not null;index" json:"user_id"` // From auth service
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Latitude  float64   `gorm:"type:decimal(10,8);not null" json:"latitude"`
	Longitude float64   `gorm:"type:decimal(11,8);not null" json:"longitude"`
	SizeAcres float64   `gorm:"type:decimal(10,2);not null" json:"size_acres"`
	SoilType  *string   `gorm:"type:varchar(50)" json:"soil_type,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Crops []Crop `gorm:"foreignKey:FieldID" json:"crops,omitempty"`
}

func (f *Field) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// Crop represents a crop planted in a field
type Crop struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FieldID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"field_id"`
	CropType        string     `gorm:"type:varchar(50);not null" json:"crop_type"`
	Variety         *string    `gorm:"type:varchar(100)" json:"variety,omitempty"`
	PlantingDate    time.Time  `gorm:"type:date;not null" json:"planting_date"`
	ExpectedHarvest *time.Time `gorm:"type:date" json:"expected_harvest,omitempty"`
	ActualHarvest   *time.Time `gorm:"type:date" json:"actual_harvest,omitempty"`
	Status          CropStatus `gorm:"type:varchar(20);not null;default:'planted'" json:"status"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Field         Field          `gorm:"foreignKey:FieldID" json:"field,omitempty"`
	HealthRecords []HealthRecord `gorm:"foreignKey:CropID" json:"health_records,omitempty"`
	Treatments    []Treatment    `gorm:"foreignKey:CropID" json:"treatments,omitempty"`
}

type CropStatus string

const (
	StatusPlanted    CropStatus = "planted"
	StatusGrowing    CropStatus = "growing"
	StatusHarvesting CropStatus = "harvesting"
	StatusHarvested  CropStatus = "harvested"
)

func (c *Crop) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// HealthRecord represents a crop health check/disease detection
type HealthRecord struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CropID          uuid.UUID `gorm:"type:uuid;not null;index" json:"crop_id"`
	CheckDate       time.Time `gorm:"autoCreateTime" json:"check_date"`
	HealthScore     int       `gorm:"check:health_score >= 0 AND health_score <= 100" json:"health_score"`
	ImageURL        string    `gorm:"type:varchar(500);not null" json:"image_url"`
	DiseaseDetected *string   `gorm:"type:varchar(100)" json:"disease_detected,omitempty"`
	Confidence      *float64  `gorm:"type:decimal(5,4)" json:"confidence,omitempty"`
	Severity        *string   `gorm:"type:varchar(20)" json:"severity,omitempty"` // mild, moderate, severe
	Notes           *string   `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	Crop       Crop        `gorm:"foreignKey:CropID" json:"crop,omitempty"`
	Treatments []Treatment `gorm:"foreignKey:HealthRecordID" json:"treatments,omitempty"`
}

func (h *HealthRecord) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}

// Treatment represents a treatment applied to a crop
type Treatment struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CropID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"crop_id"`
	HealthRecordID  *uuid.UUID `gorm:"type:uuid;index" json:"health_record_id,omitempty"`
	DiseaseName     string     `gorm:"type:varchar(100);not null" json:"disease_name"`
	TreatmentName   string     `gorm:"type:varchar(200);not null" json:"treatment_name"`
	TreatmentType   string     `gorm:"type:varchar(20);not null" json:"treatment_type"` // chemical, organic
	ApplicationDate time.Time  `gorm:"type:date;not null" json:"application_date"`
	Cost            *float64   `gorm:"type:decimal(12,2)" json:"cost,omitempty"`
	Effectiveness   *int       `gorm:"check:effectiveness IS NULL OR (effectiveness >= 1 AND effectiveness <= 5)" json:"effectiveness,omitempty"`
	Notes           *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	Crop         Crop          `gorm:"foreignKey:CropID" json:"crop,omitempty"`
	HealthRecord *HealthRecord `gorm:"foreignKey:HealthRecordID" json:"health_record,omitempty"`
}

func (t *Treatment) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
