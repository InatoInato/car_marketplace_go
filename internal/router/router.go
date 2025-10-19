package router

import (
	"github.com/InatoInato/car_marketplace_go.git/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, userHandler *handler.UserHandler, carHandler *handler.CarHandler) {
	api := app.Group("/api")

	api.Post("/register", userHandler.Register)
	api.Post("/login", userHandler.Login)
	api.Get("/users", userHandler.GetAllUsers)

	cars := api.Group("/cars")
	cars.Post("/create", carHandler.CreateCar)
	cars.Get("/all", carHandler.GetAllCars)
	cars.Get("/:id", carHandler.GetCarByID)
	cars.Put("/update/:id", carHandler.UpdateCar)
	cars.Delete("/:id", carHandler.DeleteCar)
}
