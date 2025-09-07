package app

import (
	"github.com/horlerdipo/todo-golang/internal/auth"
	"gorm.io/gorm"
)

type Container struct {
	db            *gorm.DB
	AuthContainer *auth.Container
}

func NewAppContainer(db *gorm.DB) *Container {
	return &Container{
		db:            db,
		AuthContainer: auth.NewContainer(db),
	}
}
