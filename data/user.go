package data

import (
	"context"
	"errors"
	"fmt"

	"github.com/flrn000/pc-partpicker/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore struct {
	dbPool *pgxpool.Pool
}

func NewUserStore(dbPool *pgxpool.Pool) *UserStore {
	return &UserStore{dbPool: dbPool}
}

func (us *UserStore) GetByEmail(email string) (*types.User, error) {
	result := &types.User{}

	err := us.dbPool.QueryRow(context.Background(), "SELECT * FROM users WHERE email=$1", email).Scan(
		&result.ID,
		&result.CreatedAt,
		&result.UserName,
		&result.Email,
		&result.Password,
	)
	if err != nil {
		return nil, fmt.Errorf("QueryRow failed: %v", err)
	}

	if result.ID == 0 {
		return nil, errors.New("user not found")
	}

	return result, nil
}

func (us *UserStore) Create(user *types.User) error {
	query := `
		INSERT INTO users (username, email, hashed_password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := us.dbPool.QueryRow(context.Background(), query, user.UserName, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("QueryRow failed: %v", err)
	}
	return nil
}

func (us *UserStore) Get(id int) (*types.User, error) {
	result := &types.User{}

	err := us.dbPool.QueryRow(context.Background(), "SELECT * FROM users WHERE id=$1", id).Scan(
		&result.ID,
		&result.CreatedAt,
		&result.UserName,
		&result.Email,
		&result.Password,
	)
	if err != nil {
		return nil, fmt.Errorf("QueryRow failed: %v", err)
	}

	if result.ID == 0 {
		return nil, errors.New("user not found")
	}

	return result, nil
}
