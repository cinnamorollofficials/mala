package models

import (
	"time"

	"github.com/google/uuid"
)

type VirtualKey struct {
	ID          uuid.UUID  `json:"id"`
	KeyValue    string     `json:"key_value"`
	Name        string     `json:"name"`
	TotalBudget float64    `json:"total_budget"`
	SpentAmount float64    `json:"spent_amount"`
	ExpiresAt   *time.Time `json:"expires_at"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
