package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sudeeya/avito-assignment/internal/model"
	"github.com/sudeeya/avito-assignment/internal/repository"
)

var _ Product = (*ProductService)(nil)

type ProductService struct {
	repo repository.ProductRepository
}

func newProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

// AddProduct implements Product.
func (p *ProductService) AddProduct(ctx context.Context, pvzID uuid.UUID, productType string) (model.Product, error) {
	product, err := p.repo.AddProduct(ctx, pvzID, productType)
	if errors.Is(err, repository.ErrUnsupportedProductType) {
		return model.Product{}, fmt.Errorf("adding product: %w", ErrUnsupportedProductType)
	} else if err != nil {
		return model.Product{}, ErrCannotAddProduct
	}

	return product, nil
}

// DeleteLastProduct implements Product.
func (p *ProductService) DeleteLastProduct(ctx context.Context, pvzID uuid.UUID) error {
	err := p.repo.DeleteLastProduct(ctx, pvzID)
	if errors.Is(err, repository.ErrReceptionIsEmpty) {
		return fmt.Errorf("deleting product: %w", ErrReceptionIsEmpty)
	} else if err != nil {
		return ErrCannotDeleteProduct
	}

	return nil
}
