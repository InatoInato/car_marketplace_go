package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/InatoInato/car_marketplace_go.git/internal/model"
	"github.com/InatoInato/car_marketplace_go.git/internal/repository"
	"github.com/redis/go-redis/v9"
)

type CarService struct {
	repo  *repository.CarRepository
	userRepo * repository.UserRepository
	cache *redis.Client
}

func NewCarService(repo *repository.CarRepository, userRepo *repository.UserRepository ,cache *redis.Client) *CarService {
	return &CarService{repo: repo, userRepo: userRepo, cache: cache}
}

var ctx = context.Background()


var ErrUserNotFound = errors.New("user not found")

func (s *CarService) Create(car *model.Car) error {
	user, err := s.userRepo.GetUserByID(car.UserID)
	if err != nil {
		return err // db error
	}
	if user == nil {
		return fmt.Errorf("%w: ID %d", ErrUserNotFound, car.UserID)
	}

	if err := s.repo.CreateCar(car); err != nil {
		return err
	}
	
	// Async cache invalidation
	// Very important 
	go func(id uint) {
		_ = s.cache.Del(ctx, "cars:all")
		_ = s.cache.Del(ctx, fmt.Sprintf("car:%d", id))
	}(car.ID)

	return nil
}

func (s *CarService) GetAll() ([]model.Car, error) {
	carsChan := make(chan []model.Car, 2)
	errChan := make(chan error, 2)

	// Try to get from cache
	go func(){
		val, err := s.cache.Get(ctx, "cars:all").Result()
		if err == nil{
			var cars []model.Car
			_ = json.Unmarshal([]byte(val), &cars)
			carsChan <- cars
			return
		}
		errChan <- err
		carsChan <- nil
	}()

	go func(){
		cars, err := s.repo.GetAllCars()
		if err != nil {
			errChan <-  err
			carsChan <- nil
			return
		}
		carsChan <- cars
		data, _ := json.Marshal(cars)
		_ = s.cache.Set(ctx, "cars:all", data, 5*time.Minute)
	}()

	for i := 0; i < 2; i ++{
		select{
		case cars := <- carsChan:
			if cars != nil{
				return cars, nil
			}
		case err := <- errChan:
			log.Println("Get all error: ", err)
		}
	}

	return nil, errors.New("failed to get cars from both cache and database")
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

	go func(id uint){
		s.cache.Del(ctx, "cars:all")
		s.cache.Del(ctx, fmt.Sprintf("car:%d", car.ID))
	}(car.ID)
	return nil
}

func (s *CarService) Delete(id uint) error {
	if err := s.repo.DeleteCar(id); err != nil {
		return err
	}

	go func(id uint){
		s.cache.Del(ctx, "cars:all")
		s.cache.Del(ctx, fmt.Sprintf("car:%d", id))
	}(id)
	return nil
}
