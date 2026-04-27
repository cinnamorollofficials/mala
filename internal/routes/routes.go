package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hadigunawan/mala/internal/handlers"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/health", handlers.HealthCheck)
}
