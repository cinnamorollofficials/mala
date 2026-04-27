package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hadigunawan/mala/pkg/database"
)

func ChatCompletions(c *fiber.Ctx) error {
	startTime := time.Now()
	
	// Get data from middleware Locals
	virtualKeyID := c.Locals("virtual_key_id").(string)
	apiKey := c.Locals("real_api_key").(string)
	baseURL := c.Locals("base_url").(string)
	modelDBID := c.Locals("model_db_id").(string)
	inputPrice := c.Locals("input_price").(float64)
	outputPrice := c.Locals("output_price").(float64)

	// Proxy the request
	proxyURL := baseURL + "/chat/completions"
	
	// Use c.Body() directly as it was already scrubbed by PIIScrubber
	req, _ := http.NewRequest("POST", proxyURL, bytes.NewBuffer(c.Body()))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "Failed to reach LLM provider"})
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Record Usage (Async)
	go func() {
		var openAIResp struct {
			Usage struct {
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
			} `json:"usage"`
		}
		json.Unmarshal(respBody, &openAIResp)

		cost := (float64(openAIResp.Usage.PromptTokens) * inputPrice / 1000) + 
		        (float64(openAIResp.Usage.CompletionTokens) * outputPrice / 1000)

		latency := int(time.Since(startTime).Milliseconds())

		ctx := context.Background()
		tx, _ := database.DB.Begin(ctx)
		defer tx.Rollback(ctx)

		tx.Exec(ctx,
			`INSERT INTO usage_logs (virtual_key_id, model_id, prompt_tokens, completion_tokens, total_cost, latency_ms, status_code) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			virtualKeyID, modelDBID, openAIResp.Usage.PromptTokens, openAIResp.Usage.CompletionTokens, cost, latency, resp.StatusCode,
		)

		tx.Exec(ctx, "UPDATE virtual_keys SET spent_amount = spent_amount + $1 WHERE id = $2", cost, virtualKeyID)
		
		tx.Commit(ctx)
	}()

	c.Status(resp.StatusCode)
	return c.Send(respBody)
}
