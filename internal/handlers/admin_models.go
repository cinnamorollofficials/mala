package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hadigunawan/mala/internal/models"
	"github.com/hadigunawan/mala/pkg/database"
)

func CreateModel(c *fiber.Ctx) error {
	type Request struct {
		ProviderID       string  `json:"provider_id"`
		ModelName        string  `json:"model_name"`
		InputPricePer1k  float64 `json:"input_price_per_1k"`
		OutputPricePer1k float64 `json:"output_price_per_1k"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	id := uuid.New()
	_, err := database.DB.Exec(context.Background(),
		`INSERT INTO models (id, provider_id, model_name, input_price_per_1k, output_price_per_1k) 
		 VALUES ($1, $2, $3, $4, $5)`,
		id, req.ProviderID, req.ModelName, req.InputPricePer1k, req.OutputPricePer1k,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":         id,
		"model_name": req.ModelName,
	})
}

func AdminListModels(c *fiber.Ctx) error {
	rows, err := database.DB.Query(context.Background(),
		`SELECT m.id, m.provider_id, m.model_name, m.input_price_per_1k, m.output_price_per_1k, m.is_active, m.created_at, m.updated_at, p.name as provider_name 
		 FROM models m 
		 JOIN providers p ON m.provider_id = p.id`)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	type ModelDetail struct {
		models.Model
		ProviderName string `json:"provider_name"`
	}

	modelsList := []ModelDetail{}
	for rows.Next() {
		var m ModelDetail
		err := rows.Scan(&m.ID, &m.ProviderID, &m.ModelName, &m.InputPricePer1k, &m.OutputPricePer1k, &m.IsActive, &m.CreatedAt, &m.UpdatedAt, &m.ProviderName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		modelsList = append(modelsList, m)
	}

	return c.JSON(modelsList)
}

func UpdateModel(c *fiber.Ctx) error {
	id := c.Params("id")
	type Request struct {
		InputPricePer1k  float64 `json:"input_price_per_1k"`
		OutputPricePer1k float64 `json:"output_price_per_1k"`
		IsActive         bool    `json:"is_active"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := database.DB.Exec(context.Background(),
		`UPDATE models SET input_price_per_1k = $1, output_price_per_1k = $2, is_active = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4`,
		req.InputPricePer1k, req.OutputPricePer1k, req.IsActive, id,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Model updated successfully"})
}

func DeleteModel(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := database.DB.Exec(context.Background(),
		"UPDATE models SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1",
		id,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Model deactivated successfully"})
}
