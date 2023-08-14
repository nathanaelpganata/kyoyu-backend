package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	Admin  UserRole = "admin"
	Member UserRole = "member"
)

type User struct {
	UserID    uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4();" json:"userId"`
	Email     string         `gorm:"unique;" json:"email"`
	Password  string         `gorm:"not null;" json:"password"`
	Role      UserRole       `gorm:"default:member;" json:"role"`
	Posts     []Post         `gorm:"foreignKey:AuthorID;" json:"posts"`
	CreatedAt time.Time      `gorm:"autoCreateTime;" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index;" json:"deletedAt"`
}
