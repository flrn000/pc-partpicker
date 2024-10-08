package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/flrn000/pc-partpicker/types"
	"github.com/flrn000/pc-partpicker/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenStore struct {
	dbPool *pgxpool.Pool
}

func NewRefreshTokenStore(dbPool *pgxpool.Pool) *RefreshTokenStore {
	return &RefreshTokenStore{
		dbPool: dbPool,
	}
}

func (r *RefreshTokenStore) Create(userID int, expiresAt time.Time) (types.RefreshToken, error) {
	var result types.RefreshToken

	query := `
		INSERT INTO refresh_tokens (token, user_id, expires_at)
		VALUES ($1, $2, $3)
		RETURNING created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	token, err := utils.GenerateRefreshToken()
	if err != nil {
		return result, err
	}

	err = r.dbPool.QueryRow(ctx, query, token, userID, expiresAt).Scan(&result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return result, fmt.Errorf("QueryRow failed: %v", err)
	}

	result.Value = token
	result.UserID = userID
	result.ExpiresAt = expiresAt

	return result, nil
}

func (r *RefreshTokenStore) Get(token string) (types.RefreshToken, error) {
	query := `
		SELECT * FROM refresh_tokens
		WHERE token = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	refreshToken := types.RefreshToken{}

	err := r.dbPool.QueryRow(ctx, query, token).Scan(
		&refreshToken.Value,
		&refreshToken.CreatedAt,
		&refreshToken.UpdatedAt,
		&refreshToken.UserID,
		&refreshToken.ExpiresAt,
		&refreshToken.RevokedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return refreshToken, types.ErrNoRecord
		}
		return refreshToken, err
	}

	return refreshToken, nil
}
