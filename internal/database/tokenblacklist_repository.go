package database

import (
	"gorm.io/gorm"
	"time"
)

type TokenBlacklistRepository interface {
	CheckTokenExistence(token string) bool
	InsertToken(token string, ttl *time.Time) (uint, error)
}

type tokenBlacklistRepository struct {
	db *gorm.DB
}

func NewTokenBlacklistRepository(db *gorm.DB) TokenBlacklistRepository {
	return &tokenBlacklistRepository{
		db,
	}
}

func (repo *tokenBlacklistRepository) CheckTokenExistence(token string) bool {
	//check if token exists
	tokenBlacklist := &TokenBlacklist{}
	result := repo.db.Where("token = ?", token).First(&tokenBlacklist)
	if result.Error != nil {
		return false
	}
	return true
}

func (repo *tokenBlacklistRepository) InsertToken(token string, ttl *time.Time) (uint, error) {
	tokenBlacklist := &TokenBlacklist{
		Token:     token,
		ExpiresAt: ttl,
	}

	result := repo.db.Create(&tokenBlacklist)
	if result.Error != nil {
		return 0, result.Error
	}
	return tokenBlacklist.ID, nil
}
