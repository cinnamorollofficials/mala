package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/hadigunawan/mala/pkg/database"
)

func ListModels(c *fiber.Ctx) error {
	// In a real scenario, you might filter models based on virtual_key_id permissions
	rows, err := database.DB.Query(context.Background(),
		"SELECT m.model_name, p.name as provider_name FROM models m JOIN providers p ON m.provider_id = p.id WHERE m.is_active = true")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	type ModelInfo struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		OwnedBy string `json:"owned_by"`
	}

	models := []ModelInfo{}
	for rows.Next() {
		var name, provider string
		rows.Scan(&name, &provider)
		models = append(models, ModelInfo{
			ID:      name,
			Object:  "model",
			OwnedBy: provider,
		})
	}

	return c.JSON(fiber.Map{
		"object": "list",
		"data":   models,
	})
}
