package models

import (
	"time"

	"github.com/google/uuid"
)

type Provider struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	ProviderType string    `json:"provider_type"`
	APIKey       string    `json:"api_key"`
	BaseURL      string    `json:"base_url"`
	IsActive     bool      `json:"is_active"`
	Priority     int       `json:"priority"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
