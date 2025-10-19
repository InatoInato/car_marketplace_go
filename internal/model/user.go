package model

import "gorm.io/gorm"

type User struct{
	gorm.Model
	Name string `json:"name" gorm:"type:varchar(100)"`
	Email string `json:"email" gorm:"type:varchar(100);uniqueIndex"`
	Password string `json:"-" gorm:"type:varchar(100)"`
	Cars []Car `json:"cars" gorm:"foreignKey:UserID"`
}