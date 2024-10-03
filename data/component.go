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

func (cs *ComponentStore) GetMany(limit int, componentType types.ComponentType) ([]*types.Component, error) {
	query := `
		SELECT id, created_at, updated_at, name, type, manufacturer, model, price, rating, image_path
		FROM components
		WHERE lower(type) = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := cs.dbPool.Query(ctx, query, componentType, limit)
	if err != nil {
		return nil, err
	}

	results, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[types.Component])
	if err != nil {
		return nil, err
	}

	return results, nil
}
