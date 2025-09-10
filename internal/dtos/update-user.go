package dtos

import "time"

type UpdateUserDTO struct {
	FirstName           string     `json:"first_name"`
	LastName            string     `json:"last_name"`
	Email               string     `json:"email"`
	ResetToken          *string    `json:"reset_token"`
	ResetTokenExpiresAt *time.Time `json:"reset_token_expires_at"`
}
