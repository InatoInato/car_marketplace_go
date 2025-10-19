package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/InatoInato/car_marketplace_go.git/internal/model"
	"github.com/InatoInato/car_marketplace_go.git/internal/repository"
	"github.com/redis/go-redis/v9"
)

type CarService struct {
	repo  *repository.CarRepository
	cache *redis.Client
}

func NewCarService(repo *repository.CarRepository, cache *redis.Client) *CarService {
	return &CarService{repo: repo, cache: cache}
}

var ctx = context.Background()

func (s *CarService) Create(car *model.Car) error {
	if err := s.repo.CreateCar(car); err != nil {
		return err
	}
	s.cache.Del(ctx, "cars:all") 
	return nil
}

func (s *CarService) GetAll() ([]model.Car, error) {
	val, err := s.cache.Get(ctx, "cars:all").Result()
	if err == nil {
		var cars []model.Car
		_ = json.Unmarshal([]byte(val), &cars)
		return cars, nil
	}

	cars, err := s.repo.GetAllCars()
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(cars)
	s.cache.Set(ctx, "cars:all", data, 5*time.Minute)

	return cars, nil
}

func (s *CarService) GetByID(id uint) (*model.Car, error) {
	cacheKey := fmt.Sprintf("car:%d", id)
	val, err := s.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var car model.Car
		_ = json.Unmarshal([]byte(val), &car)
		return &car, nil
	}

	car, err := s.repo.GetCarByID(id)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(car)
	s.cache.Set(ctx, cacheKey, data, 10*time.Minute)
	return car, nil
}

func (s *CarService) Update(car *model.Car) error {
	if err := s.repo.UpdateCar(car); err != nil {
		return err
	}
	s.cache.Del(ctx, "cars:all")
	s.cache.Del(ctx, fmt.Sprintf("car:%d", car.ID))
	return nil
}

func (s *CarService) Delete(id uint) error {
	if err := s.repo.DeleteCar(id); err != nil {
		return err
	}
	s.cache.Del(ctx, "cars:all")
	s.cache.Del(ctx, fmt.Sprintf("car:%d", id))
	return nil
}
