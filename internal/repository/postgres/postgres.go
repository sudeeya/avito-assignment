package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	"github.com/sudeeya/avito-assignment/internal/config"
	"github.com/sudeeya/avito-assignment/internal/model"
	"github.com/sudeeya/avito-assignment/internal/repository"
)

// Reception statuses.
const (
	_inProgressStatus = "in_progress"
	_closeStatus      = "close"
)

// Pagination defaults.
const (
	_defaultLimit  = 10
	_defaultOffset = 0
)

var _ repository.Repository = (*postgres)(nil)

type postgres struct {
	pool    *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewPostgres(ctx context.Context, cfg config.DBConfig) (*postgres, error) {
	zap.L().Info("Establishing a connection to the database...")
	pool, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		return nil, fmt.Errorf("creating pool: %w", err)
	}

	zap.L().Info("Applying migrations...")
	if err := goose.SetDialect(cfg.GooseDriver); err != nil {
		return nil, fmt.Errorf("setting dialect: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)
	if err := goose.Up(db, cfg.GooseMigrationDir); err != nil {
		return nil, fmt.Errorf("applying migrations: %w", err)
	}

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &postgres{
		pool:    pool,
		builder: builder,
	}, nil
}

// CreatePVZ implements repository.Repository.
func (p *postgres) CreatePVZ(ctx context.Context, city string) (model.PVZ, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return model.PVZ{}, fmt.Errorf("initiating transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Check if the city is supported.
	query, args, err := p.builder.
		Select("id").
		From("cities").
		Where("name = ?", city).
		ToSql()
	if err != nil {
		return model.PVZ{}, fmt.Errorf("building query: %w", err)
	}

	var cityID uuid.UUID
	err = tx.QueryRow(ctx, query, args...).Scan(&cityID)
	if errors.Is(err, pgx.ErrNoRows) { // City was not found.
		return model.PVZ{}, repository.ErrUnsupportedCity
	} else if err != nil { // Some error.
		return model.PVZ{}, fmt.Errorf("selecting city: %w", err)
	}

	// City was found, so create pvz.
	query, args, err = p.builder.
		Insert("pvzs").
		Columns("city_id").
		Values(cityID).
		Suffix("RETURNING id, registration_date").
		ToSql()
	if err != nil {
		return model.PVZ{}, fmt.Errorf("building query: %w", err)
	}

	pvz := model.PVZ{
		City: city,
	}
	err = tx.QueryRow(ctx, query, args...).Scan(&pvz.ID, &pvz.RegistrationDate)
	if err != nil {
		return model.PVZ{}, fmt.Errorf("inserting pvz: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.PVZ{}, fmt.Errorf("committing transaction: %w", err)
	}

	return pvz, nil
}

// GetPVZsByRegistrationDate implements repository.Repository.
func (p *postgres) GetPVZPagination(ctx context.Context, start, end time.Time, limit, offset int) ([]model.PVZ, error) {
	if limit <= 0 || limit > _defaultLimit {
		limit = _defaultLimit
	}

	if offset < 0 {
		offset = _defaultOffset
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("initiating transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Select pvzs with in progress receptions.
	query, args, err := p.builder.
		Select(
			"p.id",
			"p.registration_date",
			"c.name",
			"r.id AS reception_id",
			"r.datetime",
			"r.status",
		).
		Options("DISTINCT ON (p.id)").
		From("pvzs AS p").
		LeftJoin("cities AS c ON p.city_id = c.id").
		LeftJoin("receptions AS r ON r.pvz_id = p.id").
		Where("r.status = ? AND r.datetime BETWEEN ? AND ?", _inProgressStatus, start, end).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %w", err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("selecting pvzs and in progress receptions: %w", err)
	}
	defer rows.Close()

	pvzs := make([]model.PVZ, 0)
	for rows.Next() {
		var (
			pvz       model.PVZ
			reception model.Reception
		)

		err := rows.Scan(
			&pvz.ID,
			&pvz.RegistrationDate,
			&pvz.City,
			&reception.ID,
			&reception.Datetime,
			&reception.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		pvz.Receptions = append(pvz.Receptions, reception)

		pvzs = append(pvzs, pvz)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	// For each reception select its products.
	for _, pvz := range pvzs {
		query, args, err = p.builder.
			Select(
				"p.id",
				"p.datetime",
				"t.name",
			).
			From("products AS p").
			LeftJoin("product_types AS t ON p.product_type_id = t.id").
			Where("p.reception_id = ?", pvz.Receptions[0].ID).
			ToSql()
		if err != nil {
			return nil, fmt.Errorf("building query: %w", err)
		}

		rows, err := tx.Query(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("selecting products: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var product model.Product

			err := rows.Scan(
				&product.ID,
				&product.Datetime,
				&product.Type,
			)
			if err != nil {
				return nil, fmt.Errorf("scanning row: %w", err)
			}

			pvz.Receptions[0].Products = append(pvz.Receptions[0].Products, product)
		}

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("iterating rows: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	return pvzs, nil
}

// GetPVZList implements repository.Repository.
func (p *postgres) GetPVZList(ctx context.Context) ([]model.PVZ, error) {
	query, args, err := p.builder.
		Select(
			"p.id",
			"p.registration_date",
			"c.name",
		).
		From("pvzs AS p").
		LeftJoin("cities AS c ON p.city_id = c.id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %w", err)
	}

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("selecting pvzs: %w", err)
	}
	defer rows.Close()

	pvzs := make([]model.PVZ, 0)
	for rows.Next() {
		var pvz model.PVZ

		err := rows.Scan(
			&pvz.ID,
			&pvz.RegistrationDate,
			&pvz.City,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		pvzs = append(pvzs, pvz)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return pvzs, nil
}

// CreateReception implements repository.Repository.
func (p *postgres) CreateReception(ctx context.Context, pvzID uuid.UUID) (model.Reception, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return model.Reception{}, fmt.Errorf("initiating transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Check if there is a reception with "in_progress" status.
	query, args, err := p.builder.
		Select("receptions.id").
		From("receptions").
		Join("pvzs ON receptions.pvz_id = pvzs.id").
		Where("pvz_id = ? AND status = ?", pvzID, _inProgressStatus).
		ToSql()
	if err != nil {
		return model.Reception{}, fmt.Errorf("building query: %w", err)
	}

	var receptionID uuid.UUID
	err = tx.QueryRow(ctx, query, args...).Scan(&receptionID)
	if err == nil { // "in_progress" reception was found.
		return model.Reception{}, repository.ErrReceptionInProgress
	} else if !errors.Is(err, pgx.ErrNoRows) { // Error is different from ErrNoRows.
		return model.Reception{}, fmt.Errorf("selecting reception: %w", err)
	}

	// "in_progress" reception was not found, so create reception.
	query, args, err = p.builder.
		Insert("receptions").
		Columns("pvz_id").
		Values(pvzID).
		Suffix("RETURNING id, datetime, status").
		ToSql()
	if err != nil {
		return model.Reception{}, fmt.Errorf("building query: %w", err)
	}

	reception := model.Reception{
		PVZID: pvzID,
	}
	err = tx.QueryRow(ctx, query, args...).Scan(&reception.ID, &reception.Datetime, &reception.Status)
	if err != nil {
		return model.Reception{}, fmt.Errorf("inserting reception: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.Reception{}, fmt.Errorf("committing transaction: %w", err)
	}

	return reception, nil
}

// CloseLastReception implements repository.Repository.
func (p *postgres) CloseLastReception(ctx context.Context, pvzID uuid.UUID) (model.Reception, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return model.Reception{}, fmt.Errorf("initiating transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Check if there is a reception with "in_progress" status.
	query, args, err := p.builder.
		Select("receptions.id").
		From("receptions").
		Join("pvzs ON receptions.pvz_id = pvzs.id").
		Where("pvz_id = ? AND status = ?", pvzID, _inProgressStatus).
		ToSql()
	if err != nil {
		return model.Reception{}, fmt.Errorf("building query: %w", err)
	}

	var receptionID uuid.UUID
	err = tx.QueryRow(ctx, query, args...).Scan(&receptionID)
	if errors.Is(err, pgx.ErrNoRows) { // "in_progress" reception was not found.
		return model.Reception{}, repository.ErrNoReceptionInProgress
	} else if err != nil { // Some error.
		return model.Reception{}, fmt.Errorf("selecting reception: %w", err)
	}

	// "in_progress" reception was found, so close reception.
	query, args, err = p.builder.
		Update("receptions").
		Set("status", _closeStatus).
		Where("id = ?", receptionID).
		Suffix("RETURNING id, pvz_id, datetime, status").
		ToSql()
	if err != nil {
		return model.Reception{}, fmt.Errorf("building query: %w", err)
	}

	var reception model.Reception
	err = tx.QueryRow(ctx, query, args...).Scan(
		&reception.ID,
		&reception.PVZID,
		&reception.Datetime,
		&reception.Status,
	)
	if err != nil {
		return model.Reception{}, fmt.Errorf("updating reception: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.Reception{}, fmt.Errorf("committing transaction: %w", err)
	}

	return reception, nil
}

// AddProduct implements repository.Repository.
func (p *postgres) AddProduct(ctx context.Context, pvzID uuid.UUID, productType string) (model.Product, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return model.Product{}, fmt.Errorf("initiating transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Check if there is a reception with "in_progress" status.
	query, args, err := p.builder.
		Select("receptions.id").
		From("receptions").
		Join("pvzs ON receptions.pvz_id = pvzs.id").
		Where("pvz_id = ? AND status = ?", pvzID, _inProgressStatus).
		ToSql()
	if err != nil {
		return model.Product{}, fmt.Errorf("building query: %w", err)
	}

	var receptionID uuid.UUID
	err = tx.QueryRow(ctx, query, args...).Scan(&receptionID)
	if errors.Is(err, pgx.ErrNoRows) { // "in_progress" reception was not found.
		return model.Product{}, repository.ErrNoReceptionInProgress
	} else if err != nil { // Some error.
		return model.Product{}, fmt.Errorf("selecting reception: %w", err)
	}

	// "in_progress" reception was found.
	// Check if the product type is supported.
	query, args, err = p.builder.
		Select("id").
		From("product_types").
		Where("name = ?", productType).
		ToSql()
	if err != nil {
		return model.Product{}, fmt.Errorf("building query: %w", err)
	}

	var productTypeID uuid.UUID
	err = tx.QueryRow(ctx, query, args...).Scan(&productTypeID)
	if errors.Is(err, pgx.ErrNoRows) { // Product type was not found.
		return model.Product{}, repository.ErrUnsupportedProductType
	} else if err != nil { // Some error.
		return model.Product{}, fmt.Errorf("selecting city: %w", err)
	}

	// Product type was found, so add product.
	query, args, err = p.builder.
		Insert("products").
		Columns("reception_id", "product_type_id").
		Values(receptionID, productTypeID).
		Suffix("RETURNING id, datetime").
		ToSql()
	if err != nil {
		return model.Product{}, fmt.Errorf("building query: %w", err)
	}

	product := model.Product{
		ReceptionID: receptionID,
		Type:        productType,
	}
	err = tx.QueryRow(ctx, query, args...).Scan(
		&product.ID,
		&product.Datetime,
	)
	if err != nil {
		return model.Product{}, fmt.Errorf("inserting product: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.Product{}, fmt.Errorf("committing transaction: %w", err)
	}

	return product, nil
}

// DeleteLastProduct implements repository.Repository.
func (p *postgres) DeleteLastProduct(ctx context.Context, pvzID uuid.UUID) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("initiating transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Check if there is a reception with "in_progress" status.
	query, args, err := p.builder.
		Select("receptions.id").
		From("receptions").
		Join("pvzs ON receptions.pvz_id = pvzs.id").
		Where("pvz_id = ? AND status = ?", pvzID, _inProgressStatus).
		ToSql()
	if err != nil {
		return fmt.Errorf("building query: %w", err)
	}

	var receptionID uuid.UUID
	err = tx.QueryRow(ctx, query, args...).Scan(&receptionID)
	if errors.Is(err, pgx.ErrNoRows) { // "in_progress" reception was not found.
		return repository.ErrNoReceptionInProgress
	} else if err != nil { // Some error.
		return fmt.Errorf("selecting reception: %w", err)
	}

	// "in_progress" reception was found, so delete last product.
	subQuery, subArgs, err := p.builder.
		Select("id").
		From("products").
		OrderBy("datetime DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return fmt.Errorf("building subquery: %w", err)
	}

	query, args, err = p.builder.
		Delete("products").
		Where("id IN ("+subQuery+")", subArgs...).
		ToSql()
	if err != nil {
		return fmt.Errorf("building query: %w", err)
	}

	deleteTag, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("deleting product: %w", err)
	}

	if deleteTag.RowsAffected() == 0 {
		return repository.ErrReceptionIsEmpty
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}
