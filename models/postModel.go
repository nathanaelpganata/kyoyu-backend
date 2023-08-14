package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	PostID      uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4();" json:"postId"`
	Title       string         `gorm:"not null;" json:"title"`
	Slug        string         `json:"slug"`
	Body        string         `gorm:"not null;" json:"body"`
	AuthorID    uuid.UUID      `gorm:"type:uuid;" json:"authorId"`
	AuthorEmail string         `json:"authorEmail"`
	CreatedAt   time.Time      `gorm:"autoCreateTime;" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime;" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index;" json:"deletedAt"`
}
