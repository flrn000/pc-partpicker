package types

import "time"

type User struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserName  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
}
