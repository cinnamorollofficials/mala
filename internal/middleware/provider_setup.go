package middleware

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/hadigunawan/mala/pkg/database"
	"github.com/hadigunawan/mala/pkg/utils"
)

func ProviderSetup(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodPost {
		return c.Next()
	}

	var requestBody map[string]interface{}
	if err := json.Unmarshal(c.Body(), &requestBody); err != nil {
		return c.Next() // Let the handler handle bad body
	}

	modelName, ok := requestBody["model"].(string)
	if !ok {
		return c.Next()
	}

	var providerID, providerType, encryptedKey, baseURL string
	var modelID string
	var inputPrice, outputPrice float64

	err := database.DB.QueryRow(context.Background(),
		`SELECT p.id, p.provider_type, p.api_key, p.base_url, m.id, m.input_price_per_1k, m.output_price_per_1k 
		 FROM models m 
		 JOIN providers p ON m.provider_id = p.id 
		 WHERE m.model_name = $1 AND m.is_active = true AND p.is_active = true`,
		modelName,
	).Scan(&providerID, &providerType, &encryptedKey, &baseURL, &modelID, &inputPrice, &outputPrice)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Model not found or provider inactive"})
	}

	apiKey, err := utils.Decrypt(encryptedKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal security error"})
	}

	// Store for handler
	c.Locals("real_api_key", apiKey)
	c.Locals("base_url", baseURL)
	c.Locals("provider_type", providerType)
	c.Locals("model_db_id", modelID)
	c.Locals("input_price", inputPrice)
	c.Locals("output_price", outputPrice)

	return c.Next()
}
