package handler

import (
	"github.com/InatoInato/car_marketplace_go.git/internal/config"
	"github.com/InatoInato/car_marketplace_go.git/internal/dto"
	"github.com/InatoInato/car_marketplace_go.git/internal/model"
	"github.com/InatoInato/car_marketplace_go.git/internal/service"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct{
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(ctx *fiber.Ctx) error {
	var input dto.RegisterInput

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	if err := config.ValidateStruct(&input); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Validation failed",
		})
	}

	user := &model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := h.service.Register(user); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return ctx.Status(201).JSON(fiber.Map{
		"message": "User registered successfully",
	})
}

func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	var input struct{
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	token, err := h.service.Login(input.Email, input.Password)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	return ctx.JSON(fiber.Map{
		"token": token,
	})
}