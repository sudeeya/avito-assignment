package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sudeeya/avito-assignment/internal/config"
	"github.com/sudeeya/avito-assignment/internal/model"
	"github.com/sudeeya/avito-assignment/internal/repository"
)

type Auth interface {
	IssueToken(ctx context.Context) (string, error)
	VerifyToken(ctx context.Context, token string) error
}

type PVZ interface {
	CreatePVZ(ctx context.Context, city string) (model.PVZ, error)
	GetPVZPagination(ctx context.Context, start, end time.Time, limit, offset int) ([]model.PVZ, error)
}

type Reception interface {
	CreateReception(ctx context.Context, pvzID uuid.UUID) (model.Reception, error)
	CloseLastReception(ctx context.Context, pvzID uuid.UUID) (model.Reception, error)
}

type Product interface {
	AddProduct(ctx context.Context, pvzID uuid.UUID, productType string) (model.Product, error)
	DeleteLastProduct(ctx context.Context, pvzID uuid.UUID) error
}

type Services struct {
	Auth      Auth
	PVZ       PVZ
	Reception Reception
	Product   Product
}

func NewService(cfg config.ServerConfig, repo repository.Repository) (*Services, error) {
	auth, err := newAuthService(cfg)
	if err != nil {
		return nil, err
	}

	return &Services{
		Auth:      auth,
		PVZ:       newPVZService(repo),
		Reception: newReceptionService(repo),
		Product:   newProductService(repo),
	}, nil
}
