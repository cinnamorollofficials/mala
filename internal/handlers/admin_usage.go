package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/hadigunawan/mala/internal/models"
	"github.com/hadigunawan/mala/pkg/database"
)

func GetUsageSummary(c *fiber.Ctx) error {
	var totalCost float64
	err := database.DB.QueryRow(context.Background(), 
		"SELECT COALESCE(SUM(total_cost), 0) FROM usage_logs WHERE created_at >= CURRENT_DATE").Scan(&totalCost)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"total_cost_today": totalCost,
	})
}

func GetKeyUsageHistory(c *fiber.Ctx) error {
	keyID := c.Params("key_id")
	
	rows, err := database.DB.Query(context.Background(),
		`SELECT id, virtual_key_id, model_id, prompt_tokens, completion_tokens, total_cost, latency_ms, status_code, created_at 
		 FROM usage_logs 
		 WHERE virtual_key_id = $1 
		 ORDER BY created_at DESC 
		 LIMIT 100`,
		keyID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	logs := []models.UsageLog{}
	for rows.Next() {
		var l models.UsageLog
		err := rows.Scan(&l.ID, &l.VirtualKeyID, &l.ModelID, &l.PromptTokens, &l.CompletionTokens, &l.TotalCost, &l.LatencyMS, &l.StatusCode, &l.CreatedAt)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		logs = append(logs, l)
	}

	return c.JSON(logs)
}
