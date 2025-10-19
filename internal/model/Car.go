package model

import "gorm.io/gorm"

type Car struct{
	gorm.Model
	Make string `json:"make"`
	CarModel string `json:"car_model"`
	Year int `json:"year"`
	Color string `json:"color"`
	EngineCapacity float64 `json:"engine_capacity"`
	Transmission string `json:"transmission"`
	IsRunning bool `json:"is_running"`
	UserID uint `json:"user_id" gorm:"index"`
}
