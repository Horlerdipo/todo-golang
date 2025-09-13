package database

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	FirstName           string     `json:"first_name"`
	LastName            string     `json:"last_name"`
	Email               string     `json:"email"`
	Password            string     `json:"-"`
	ResetToken          *string    `json:"reset_token"`
	ResetTokenExpiresAt *time.Time `json:"reset_token_expires_at"`
}
