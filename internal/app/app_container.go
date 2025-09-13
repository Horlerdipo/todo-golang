package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/horlerdipo/todo-golang/internal/auth"
	"github.com/horlerdipo/todo-golang/internal/todo"
	"gorm.io/gorm"
)

type Container struct {
	db            *gorm.DB
	AuthContainer *auth.Container
	TodoContainer *todo.Container
}

func NewAppContainer(db *gorm.DB) *Container {
	return &Container{
		db:            db,
		AuthContainer: auth.NewContainer(db),
		TodoContainer: todo.NewContainer(db),
	}
}

func (container *Container) RegisterRoutes(r *chi.Mux) {
	container.AuthContainer.RegisterRoutes(r)
	container.TodoContainer.RegisterRoutes(r)
}
