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

func (r *CarRepository) CreateCarPhoto(photo *model.CarPhoto) error {
	return r.db.Create(photo).Error
}

func (r *CarRepository) GetCarPhotos(carID uint) ([]model.CarPhoto, error) {
	var photos []model.CarPhoto
	err := r.db.Where("car_id = ?", carID).Find(&photos).Error
	if err != nil {
		return nil, err
	}
	return photos, nil
}

func (r *CarRepository) GetCarByID(id uint) (*model.Car, error) {
	var car model.Car
	err := r.db.First(&car, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &car, nil
}

func (r *CarRepository) GetAllCars() ([]model.Car, error) {
	var cars []model.Car
	err := r.db.Find(&cars).Error
	if err != nil {
		return nil, err
	}
	return cars, nil
}

func (r *CarRepository) UpdateCar(car *model.Car) error {
	return r.db.Save(car).Error
}

func (r *CarRepository) DeleteCar(id uint) error {
	result := r.db.Delete(&model.Car{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		// Авто не найдено
		return nil
	}
	return nil
}
