package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `json:"id"`
	ReceptionID uuid.UUID `json:"reception_id,omitempty,omitzero"`
	Datetime    time.Time `json:"datetime"`
	Type        string    `json:"type"`
}
