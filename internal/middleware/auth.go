package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/hadigunawan/mala/pkg/database"
)

func VirtualKeyMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing Authorization header",
		})
	}

	// Format: Bearer sk-gh-...
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid Authorization header format",
		})
	}

	apiKey := parts[1]

	var id string
	var isActive bool
	var spentAmount float64
	var totalBudget float64

	err := database.DB.QueryRow(context.Background(),
		"SELECT id, is_active, spent_amount, total_budget FROM virtual_keys WHERE key_value = $1",
		apiKey,
	).Scan(&id, &isActive, &spentAmount, &totalBudget)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid Virtual Key",
		})
	}

	if !isActive {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Virtual Key is inactive",
		})
	}

	if totalBudget > 0 && spentAmount >= totalBudget {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Budget exceeded for this Virtual Key",
		})
	}

	// Store virtual key info in context for later use
	c.Locals("virtual_key_id", id)

	return c.Next()
}
