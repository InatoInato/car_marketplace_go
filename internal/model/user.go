package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint `json:"user_id" gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name     string `gorm:"type:varchar(100);not null"`
	Email    string `gorm:"type:varchar(100);unique;not null"`
	Password string `gorm:"type:varchar(100);not null"`
	Role     string `gorm:"type:text;default:'user'"`
	PhotoURL string `gorm:"type:text"`
}
