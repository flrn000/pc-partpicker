package types

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"
)

type AppConfig struct {
	Logger    *slog.Logger
	JWTSecret string
	Env       string
	DB_URL    string
	Address   string
}

type ComponentType string

const (
	TypeDefault     ComponentType = ""
	TypeCPU         ComponentType = "procesor"
	TypeGPU         ComponentType = "placa-video"
	TypeMotherboard ComponentType = "placa-de-baza"
	TypeMemory      ComponentType = "memorie"
	TypeSSD         ComponentType = "ssd"
	TypeHDD         ComponentType = "hdd"
	TypeCPUCooler   ComponentType = "cooler"
	TypePSU         ComponentType = "sursa"
	TypeCase        ComponentType = "carcasa"
)

type ContextKey string

const UserContextKey = ContextKey("user")

var (
	ErrNoRecord           = errors.New("no matching record found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Filters struct {
	Page     int
	PageSize int
	Query    string
	Sort     string
}

type User struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserName  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
}

type Component struct {
	ID           int       `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Manufacturer string    `json:"manufacturer"`
	Model        string    `json:"model"`
	Price        string    `json:"price"`
	Rating       int16     `json:"rating"`
	ImageURL     string    `json:"image_url"`
}

type CreateProductPayload struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	Price        string `json:"price"`
	Rating       int16  `json:"rating"`
	ImageURL     string `json:"image_url"`
}

type RefreshToken struct {
	Value     string       `json:"refresh_token,omitempty"`
	CreatedAt time.Time    `json:"created_at,omitempty"`
	UpdatedAt time.Time    `json:"updated_at,omitempty"`
	UserID    int          `json:"-"`
	ExpiresAt time.Time    `json:"expires_at,omitempty"`
	RevokedAt sql.NullTime `json:"-"`
}
