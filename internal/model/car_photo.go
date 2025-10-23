package model

import "gorm.io/gorm"

type CarPhoto struct {
	gorm.Model
	CarID uint   `json:"car_id" gorm:"index"`
	URL   string `json:"url" gorm:"not null"`
}
