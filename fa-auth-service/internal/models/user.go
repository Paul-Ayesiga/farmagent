package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleFarmer           UserRole = "farmer"
	RoleExtensionOfficer UserRole = "extension_officer"
	RoleBuyer            UserRole = "buyer"
	RoleAdmin            UserRole = "admin"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Phone        *string        `gorm:"uniqueIndex;size:20" json:"phone,omitempty"`
	Email        *string        `gorm:"uniqueIndex;size:255" json:"email,omitempty"`
	PasswordHash string         `gorm:"-" json:"-"`
	Password     string         `gorm:"column:password;not null" json:"-"`
	FirstName    string         `gorm:"size:100;not null" json:"first_name"`
	LastName     string         `gorm:"size:100;not null" json:"last_name"`
	Role         UserRole       `gorm:"type:varchar(20);default:farmer" json:"role"`
	District     *string        `gorm:"size:100" json:"district,omitempty"`
	Language     string         `gorm:"size:10;default:en" json:"language"`
	FarmSize     *float64       `gorm:"type:decimal(10,2)" json:"farm_size,omitempty"`
	IsVerified   bool           `gorm:"default:false" json:"is_verified"`
	VerifiedAt   *time.Time     `json:"verified_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (u *User) TableName() string {
	return "users"
}
