package model

import (
	"time"

	"github.com/google/uuid"
)

type Reception struct {
	ID       uuid.UUID `json:"id"`
	PVZID    uuid.UUID `json:"pvz_id,omitempty,omitzero"`
	Datetime time.Time `json:"datetime"`
	Status   string    `json:"status"`
	Products []Product `json:"products,omitempty,omitzero"`
}
