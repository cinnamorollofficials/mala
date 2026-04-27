package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/hadigunawan/mala/pkg/database"
)

func HealthCheck(c *fiber.Ctx) error {
	status := "UP"
	dbStatus := "UP"

	err := database.DB.Ping(context.Background())
	if err != nil {
		dbStatus = "DOWN"
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":   status,
		"database": dbStatus,
	})
}
