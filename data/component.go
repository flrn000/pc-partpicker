package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/flrn000/pc-partpicker/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ComponentStore struct {
	dbPool *pgxpool.Pool
}

func NewComponentStore(dbPool *pgxpool.Pool) *ComponentStore {
	return &ComponentStore{dbPool: dbPool}
}

func (cs *ComponentStore) Create(component *types.Component) error {
	query := `
		INSERT INTO components (name, type, manufacturer, model, price, rating, image_path)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := cs.dbPool.QueryRow(
		ctx,
		query,
		component.Name,
		component.Type,
		component.Manufacturer,
		component.Model,
		component.Price,
		component.Rating,
		component.ImageURL,
	).Scan(&component.ID, &component.CreatedAt, &component.UpdatedAt)

	if err != nil {
		return fmt.Errorf("QueryRow failed: %v", err)
	}

	return nil
}

func (cs *ComponentStore) Get(id int) (*types.Component, error) {
	result := &types.Component{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := cs.dbPool.QueryRow(ctx, "SELECT * FROM components WHERE id=$1", id).Scan(
		&result.ID,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.Name,
		&result.Type,
		&result.Manufacturer,
		&result.Model,
		&result.Price,
		&result.Rating,
		&result.ImageURL,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return result, nil
}

func (cs *ComponentStore) GetAll(componentType types.ComponentType, filters types.Filters) ([]*types.Component, error) {
	query := `
		SELECT *
		FROM components
		WHERE lower(type) = $1 OR $1 = ''
		AND to_tsvector('simple', name) @@ plainto_tsquery('simple', $2) OR $2 = ''
		ORDER BY $3 DESC
		LIMIT $4
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := cs.dbPool.Query(ctx, query, componentType, filters.Query, filters.Sort, filters.PageSize)
	if err != nil {
		return nil, err
	}

	results, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[types.Component])
	if err != nil {
		return nil, err
	}

	return results, nil
}
