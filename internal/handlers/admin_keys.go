package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hadigunawan/mala/internal/models"
	"github.com/hadigunawan/mala/pkg/database"
)

func CreateVirtualKey(c *fiber.Ctx) error {
	type Request struct {
		Name        string  `json:"name"`
		TotalBudget float64 `json:"total_budget"`
		ExpiresIn   int     `json:"expires_in_days"` // Days from now
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	id := uuid.New()
	keyValue := "sk-gh-" + uuid.New().String()
	
	var expiresAt *time.Time
	if req.ExpiresIn > 0 {
		t := time.Now().AddDate(0, 0, req.ExpiresIn)
		expiresAt = &t
	}

	_, err := database.DB.Exec(context.Background(),
		"INSERT INTO virtual_keys (id, key_value, name, total_budget, expires_at) VALUES ($1, $2, $3, $4, $5)",
		id, keyValue, req.Name, req.TotalBudget, expiresAt,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":        id,
		"key_value": keyValue,
		"name":      req.Name,
	})
}

func ListVirtualKeys(c *fiber.Ctx) error {
	rows, err := database.DB.Query(context.Background(), 
		"SELECT id, key_value, name, total_budget, spent_amount, is_active, expires_at FROM virtual_keys")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	keys := []models.VirtualKey{}
	for rows.Next() {
		var k models.VirtualKey
		err := rows.Scan(&k.ID, &k.KeyValue, &k.Name, &k.TotalBudget, &k.SpentAmount, &k.IsActive, &k.ExpiresAt)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		keys = append(keys, k)
	}

	return c.JSON(keys)
}

func TopupVirtualKey(c *fiber.Ctx) error {
	id := c.Params("id")
	type Request struct {
		Amount float64 `json:"amount"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := database.DB.Exec(context.Background(),
		"UPDATE virtual_keys SET total_budget = total_budget + $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
		req.Amount, id,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Topup successful"})
}

func DeleteVirtualKey(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := database.DB.Exec(context.Background(),
		"UPDATE virtual_keys SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1",
		id,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Key deactivated successfully"})
}
