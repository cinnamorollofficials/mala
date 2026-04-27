package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hadigunawan/mala/pkg/database"
	"github.com/hadigunawan/mala/pkg/utils"
)

func CreateProvider(c *fiber.Ctx) error {
	type Request struct {
		Name         string `json:"name"`
		ProviderType string `json:"provider_type"`
		APIKey       string `json:"api_key"`
		BaseURL      string `json:"base_url"`
		Priority     int    `json:"priority"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	encryptedKey, err := utils.Encrypt(req.APIKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encrypt API key"})
	}

	id := uuid.New()
	_, err = database.DB.Exec(context.Background(),
		"INSERT INTO providers (id, name, provider_type, api_key, base_url, priority) VALUES ($1, $2, $3, $4, $5, $6)",
		id, req.Name, req.ProviderType, encryptedKey, req.BaseURL, req.Priority,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":   id,
		"name": req.Name,
	})
}

func ProvidersHealth(c *fiber.Ctx) error {
	rows, err := database.DB.Query(context.Background(), "SELECT name, base_url FROM providers WHERE is_active = true")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	type HealthReport struct {
		Name   string `json:"name"`
		URL    string `json:"url"`
		Status string `json:"status"`
	}

	reports := []HealthReport{}
	client := http.Client{Timeout: 5 * time.Second}

	for rows.Next() {
		var name, url string
		rows.Scan(&name, &url)
		
		status := "UP"
		resp, err := client.Get(url)
		if err != nil || resp.StatusCode >= 500 {
			status = "DOWN"
		}
		if resp != nil {
			resp.Body.Close()
		}

		reports = append(reports, HealthReport{Name: name, URL: url, Status: status})
	}

	return c.JSON(reports)
}
