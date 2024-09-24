package types

import (
	"errors"
	"time"
)

type ComponentType string

const (
	TypeCPU         ComponentType = "procesor"
	TypeGPU         ComponentType = "placa video"
	TypeMotherboard ComponentType = "placa de baza"
	TypeMemory      ComponentType = "memorie"
	TypeSSD         ComponentType = "ssd"
	TypeHDD         ComponentType = "hdd"
	TypeCPUCooler   ComponentType = "cooler"
	TypePSU         ComponentType = "sursa"
	TypeCase        ComponentType = "carcasa"
)

var ErrNoRecord error = errors.New("no matching record found")

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
