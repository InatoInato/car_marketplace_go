package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/InatoInato/car_marketplace_go.git/internal/model"
	"github.com/InatoInato/car_marketplace_go.git/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo  *repository.UserRepository
	cache *redis.Client
}

func NewUserService(repo *repository.UserRepository, cache *redis.Client) *UserService {
	return &UserService{repo: repo, cache: cache}
}

func (s *UserService) Register(user *model.User) error {
	existingUser, err := s.repo.GetUserByEmail(user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil && existingUser.ID != 0 {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	if err := s.repo.CreateUser(user); err != nil {
		return err
	}

	go func(id uint) {
		_ = s.cache.Del(ctx, "users:all")
		_ = s.cache.Del(ctx, fmt.Sprintf("user:%d", id))
	}(user.ID)

	return nil
}

func (s *UserService) Login(email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil || user == nil {
		return "", fmt.Errorf("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	val, err := s.cache.Get(ctx, "users:all").Result()

	// Try to get from cache
	if err == nil {
		var users []model.User
		if err := json.Unmarshal([]byte(val), &users); err == nil {
			log.Println("Cache hit: users:all")
			return users, nil
		}
		log.Println("Cache corrupted, fallback to DB")
	}

	if err != redis.Nil && err != nil {
		log.Println("Cache error:", err)
	}

	users, dbErr := s.repo.GetAllUsers()
	if dbErr != nil {
		return nil, dbErr
	}

	go func(users []model.User) {
		data, _ := json.Marshal(users)
		if err := s.cache.Set(ctx, "users:all", data, 5*time.Minute).Err(); err != nil {
			log.Printf("Failed to set cache for users:all: %v", err)
		}
	}(users)

	return users, nil
}

func (s *UserService) UpdateUser(user *model.User) error {
	if err := s.repo.UpdateUser(user); err != nil {
		return err
	}

	go func(id uint) {
		_ = s.cache.Del(ctx, "users:all")
		_ = s.cache.Del(ctx, fmt.Sprintf("user:%d", id))
	}(user.ID)

	return nil
}

func (s *UserService) DeleteUser(id uint) error {
	if err := s.repo.DeleteUserById(id); err != nil {
		return err
	}

	go func(id uint) {
		_ = s.cache.Del(ctx, "users:all")
		_ = s.cache.Del(ctx, fmt.Sprintf("user:%d", id))
	}(id)

	return nil
}
