package model

import "gorm.io/gorm"

type Car struct {
	gorm.Model
	Make           string  `json:"make" validate:"required"`                    
	CarModel       string  `json:"car_model" validate:"required"`               
	Year           int     `json:"year" validate:"required,min=1900,max=2100"`  
	Color          string  `json:"color" validate:"required"`                   
	EngineCapacity float64 `json:"engine_capacity" validate:"required,gt=0"`    
	Transmission   string  `json:"transmission" validate:"required,oneof=AT MT CVT DSG"` 
	IsRunning      bool    `json:"is_running"`                                  
	Price          float64 `json:"price" validate:"required,gt=0"`              
	UserID         uint    `json:"user_id" gorm:"index" validate:"required"`    // owner ID
}
