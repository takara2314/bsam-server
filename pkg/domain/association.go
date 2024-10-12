package domain

import "time"

type Association struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	ContractType string    `json:"contract_type"`
	ExpiresAt    time.Time `json:"expires_at"`
}
