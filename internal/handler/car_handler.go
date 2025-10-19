package handler

import (
	"errors"
	"strconv"

	"github.com/InatoInato/car_marketplace_go.git/internal/model"
	"github.com/InatoInato/car_marketplace_go.git/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CarHandler struct {
	service   *service.CarService
	validator *validator.Validate
}

func NewCarHandler(service *service.CarService) *CarHandler {
	return &CarHandler{
		service:   service,
		validator: validator.New(),
	}
}

// ✅ CREATE
func (h *CarHandler) CreateCar(c *fiber.Ctx) error {
	var car model.Car
	if err := c.BodyParser(&car); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.validator.Struct(car); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	
	err := h.service.Create(&car)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(car)
}


// ✅ READ ALL
func (h *CarHandler) GetAllCars(c *fiber.Ctx) error {
	cars, err := h.service.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(cars)
}

// ✅ READ BY ID
func (h *CarHandler) GetCarByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid car ID"})
	}

	car, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Car not found"})
	}
	return c.JSON(car)
}

// ✅ UPDATE
func (h *CarHandler) UpdateCar(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid car ID"})
	}

	var car model.Car
	if err := c.BodyParser(&car); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	car.ID = uint(id)

	if err := h.validator.Struct(car); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.Update(&car); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(car)
}

// ✅ DELETE
func (h *CarHandler) DeleteCar(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid car ID"})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
