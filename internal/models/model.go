package models

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ID               uuid.UUID `json:"id"`
	ProviderID       uuid.UUID `json:"provider_id"`
	ModelName        string    `json:"model_name"`
	InputPricePer1k  float64   `json:"input_price_per_1k"`
	OutputPricePer1k float64   `json:"output_price_per_1k"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
