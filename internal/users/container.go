package users

import (
	"gorm.io/gorm"
)

type Container struct {
	UserRepository UserRepository
}

func NewContainer(db *gorm.DB) *Container {
	return &Container{
		UserRepository: NewUserRepository(db),
	}
}
