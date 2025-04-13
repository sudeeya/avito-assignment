package model

import (
	"time"

	"github.com/google/uuid"
)

type PVZ struct {
	ID               uuid.UUID   `json:"id"`
	RegistrationDate time.Time   `json:"registration_date"`
	City             string      `json:"city"`
	Receptions       []Reception `json:"receptions,omitempty,omitzero"`
}
