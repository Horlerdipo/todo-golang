package database

import (
	"time"
)

type User struct {
	Model
	FirstName           string     `json:"first_name"`
	LastName            string     `json:"last_name"`
	Email               string     `json:"email"`
	Password            string     `json:"-"`
	ResetToken          *string    `json:"reset_token"`
	ResetTokenExpiresAt *time.Time `json:"reset_token_expires_at"`
	Todos               []Todo     `json:"todos"`
}
