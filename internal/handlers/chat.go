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
	"github.com/hadigunawan/mala/pkg/utils"
)

func ChatCompletions(c *fiber.Ctx) error {
	startTime := time.Now()
	virtualKeyID := c.Locals("virtual_key_id").(string)

	var requestBody map[string]interface{}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	modelName, ok := requestBody["model"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing model field"})
	}

	// 1. Find provider and model details
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

	// 2. Decrypt API Key
	apiKey, err := utils.Decrypt(encryptedKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal security error"})
	}

	// 3. Prepare Proxy Request
	jsonBody, _ := json.Marshal(requestBody)
	proxyURL := baseURL + "/chat/completions"
	
	req, _ := http.NewRequest("POST", proxyURL, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "Failed to reach LLM provider"})
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// 4. Record Usage (Simplified)
	// In a real app, you'd parse usage from OpenAI response
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

	// Async log usage to not block response
	go func() {
		ctx := context.Background()
		tx, _ := database.DB.Begin(ctx)
		defer tx.Rollback(ctx)

		// Insert log
		tx.Exec(ctx,
			`INSERT INTO usage_logs (virtual_key_id, model_id, prompt_tokens, completion_tokens, total_cost, latency_ms, status_code) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			virtualKeyID, modelID, openAIResp.Usage.PromptTokens, openAIResp.Usage.CompletionTokens, cost, latency, resp.StatusCode,
		)

		// Update virtual key spent amount
		tx.Exec(ctx, "UPDATE virtual_keys SET spent_amount = spent_amount + $1 WHERE id = $2", cost, virtualKeyID)
		
		tx.Commit(ctx)
	}()

	// 5. Return Response
	c.Status(resp.StatusCode)
	return c.Send(respBody)
}
