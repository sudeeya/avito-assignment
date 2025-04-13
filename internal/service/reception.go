package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sudeeya/avito-assignment/internal/model"
	"github.com/sudeeya/avito-assignment/internal/repository"
)

var _ Reception = (*ReceptionService)(nil)

type ReceptionService struct {
	repo repository.ReceptionRepository
}

func newReceptionService(repo repository.ReceptionRepository) *ReceptionService {
	return &ReceptionService{
		repo: repo,
	}
}

// CloseLastReception implements Reception.
func (r *ReceptionService) CloseLastReception(ctx context.Context, pvzID uuid.UUID) (model.Reception, error) {
	reception, err := r.repo.CloseLastReception(ctx, pvzID)
	if errors.Is(err, repository.ErrNoReceptionInProgress) {
		return model.Reception{}, fmt.Errorf("closing reception: %w", ErrNoReceptionInProgress)
	} else if err != nil {
		return model.Reception{}, ErrCannotCloseReception
	}

	return reception, nil
}

// CreateReception implements Reception.
func (r *ReceptionService) CreateReception(ctx context.Context, pvzID uuid.UUID) (model.Reception, error) {
	reception, err := r.repo.CreateReception(ctx, pvzID)
	if errors.Is(err, repository.ErrReceptionInProgress) {
		return model.Reception{}, fmt.Errorf("creating reception: %w", ErrReceptionInProgress)
	} else if err != nil {
		return model.Reception{}, ErrCannotCreateReception
	}

	return reception, nil
}
