package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sudeeya/avito-assignment/internal/model"
	"github.com/sudeeya/avito-assignment/internal/repository"
)

var _ PVZ = (*PVZService)(nil)

type PVZService struct {
	repo repository.PVZRepository
}

func newPVZService(repo repository.PVZRepository) *PVZService {
	return &PVZService{
		repo: repo,
	}
}

// CreatePVZ implements PVZ.
func (p *PVZService) CreatePVZ(ctx context.Context, city string) (model.PVZ, error) {
	pvz, err := p.repo.CreatePVZ(ctx, city)
	if errors.Is(err, repository.ErrUnsupportedCity) {
		return model.PVZ{}, fmt.Errorf("creating pvz: %w", ErrUnsupportedCity)
	} else if err != nil {
		return model.PVZ{}, ErrCannotCreatePVZ
	}

	return pvz, nil
}

// GetPVZPagination implements PVZ.
func (p *PVZService) GetPVZPagination(ctx context.Context, start time.Time, end time.Time, limit int, offset int) ([]model.PVZ, error) {
	pvzs, err := p.repo.GetPVZPagination(ctx, start, end, limit, offset)
	if err != nil {
		return nil, ErrCannotGetPVZ
	}

	return pvzs, nil
}

// GetPVZList implements PVZ.
func (p *PVZService) GetPVZList(ctx context.Context) ([]model.PVZ, error) {
	pvzs, err := p.repo.GetPVZList(ctx)
	if err != nil {
		return nil, ErrCannotGetPVZ
	}

	return pvzs, nil
}
