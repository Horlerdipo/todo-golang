package database

import (
	"time"
)

type TokenBlacklist struct {
	Model
	Token     string
	ExpiresAt *time.Time //TODO: we will use to run a scheduler that runs every day and removes expired tokens from here so we don't have a large table
}
