package model

import (
	"time"

	"gorm.io/gorm"
)

// Note model - `notes` table
type Note struct {
	NoteID    uint64         `gorm:"primaryKey" json:"noteID,omitempty"`
	CreatedAt time.Time      `json:"createdAt,omitempty"`
	UpdatedAt time.Time      `json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Title     string         `json:"title,omitempty"`
	Body      string         `json:"body,omitempty"`
	IDUser    uint64         `json:"-"`
}
