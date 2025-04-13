package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/sudeeya/avito-assignment/internal/model"
)

type Repository interface {
	PVZRepository
	ReceptionRepository
	ProductRepository
}

type PVZRepository interface {
	CreatePVZ(ctx context.Context, city string) (model.PVZ, error)
	GetPVZPagination(ctx context.Context, start, end time.Time, limit, offset int) ([]model.PVZ, error)
	GetPVZList(ctx context.Context) ([]model.PVZ, error)
}

type ReceptionRepository interface {
	CreateReception(ctx context.Context, pvzID uuid.UUID) (model.Reception, error)
	CloseLastReception(ctx context.Context, pvzID uuid.UUID) (model.Reception, error)
}

type ProductRepository interface {
	AddProduct(ctx context.Context, pvzID uuid.UUID, productType string) (model.Product, error)
	DeleteLastProduct(ctx context.Context, pvzID uuid.UUID) error
}
