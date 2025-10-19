package repository

import (
	"github.com/InatoInato/car_marketplace_go.git/internal/model"
	"gorm.io/gorm"
)

type CarRepository struct {
	db *gorm.DB
}

func NewCarRepository(db *gorm.DB) *CarRepository {
	return &CarRepository{db: db}
}

func (r *CarRepository) CreateCar(car *model.Car) error {
	return r.db.Create(car).Error
}

func (r *CarRepository) GetCarByID(id uint) (*model.Car, error) {
	var car model.Car
	err := r.db.First(&car, id).Error
	return &car, err
}

func (r *CarRepository) GetAllCars() ([]model.Car, error) {
	var cars []model.Car
	err := r.db.Find(&cars).Error
	return cars, err
}

func (r *CarRepository) UpdateCar(car *model.Car) error {
	return r.db.Save(car).Error
}

func (r *CarRepository) DeleteCar(id uint) error {
	return r.db.Delete(&model.Car{}, id).Error
}
