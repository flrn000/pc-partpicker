package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPSQLStorage(dbURL string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), dbURL)

	return dbpool, err
}
