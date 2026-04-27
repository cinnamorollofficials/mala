package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
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
		"SELECT id, model_id, prompt_tokens, completion_tokens, total_cost, status_code, created_at FROM usage_logs WHERE virtual_key_id = $1 ORDER BY created_at DESC LIMIT 100",
		keyID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	type LogEntry struct {
		ID               int64   `json:"id"`
		ModelID          string  `json:"model_id"`
		PromptTokens     int     `json:"prompt_tokens"`
		CompletionTokens int     `json:"completion_tokens"`
		TotalCost        float64 `json:"total_cost"`
		StatusCode       int     `json:"status_code"`
		CreatedAt        string  `json:"created_at"`
	}

	logs := []LogEntry{}
	for rows.Next() {
		var l LogEntry
		var createdAt interface{} // Handle timestamp
		rows.Scan(&l.ID, &l.ModelID, &l.PromptTokens, &l.CompletionTokens, &l.TotalCost, &l.StatusCode, &createdAt)
		logs = append(logs, l)
	}

	return c.JSON(logs)
}
