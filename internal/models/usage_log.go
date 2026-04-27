package models

import (
	"time"

	"github.com/google/uuid"
)

type UsageLog struct {
	ID               int64     `json:"id"`
	VirtualKeyID     uuid.UUID `json:"virtual_key_id"`
	ModelID          uuid.UUID `json:"model_id"`
	PromptTokens     int       `json:"prompt_tokens"`
	CompletionTokens int       `json:"completion_tokens"`
	TotalCost        float64   `json:"total_cost"`
	LatencyMS        int       `json:"latency_ms"`
	StatusCode       int       `json:"status_code"`
	CreatedAt        time.Time `json:"created_at"`
}
