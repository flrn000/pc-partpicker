package types

import (
	"time"
)

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
	Rating       uint8     `json:"rating"`
	ImageURL     string    `json:"image_url"`
}

type CreateProductPayload struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	Price        string `json:"price"`
	Rating       uint8  `json:"rating"`
	ImageURL     string `json:"image_url"`
}
