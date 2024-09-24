package data

import (
	"context"
	"fmt"

	"github.com/flrn000/pc-partpicker/types"
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
	err := cs.dbPool.QueryRow(
		context.Background(),
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
