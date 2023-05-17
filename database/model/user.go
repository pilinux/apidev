// Package model hosts the database models
package model

import (
	"time"

	"gorm.io/gorm"
)

// User model - `users` table
type User struct {
	UserID    uint64         `gorm:"primaryKey" json:"userID,omitempty"`
	CreatedAt time.Time      `json:"createdAt,omitempty"`
	UpdatedAt time.Time      `json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	NickName  string         `json:"nickName,omitempty"`
	IDAuth    uint64         `json:"-"`
	Notes     []Note         `gorm:"foreignkey:IDUser;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"notes,omitempty"`
}
