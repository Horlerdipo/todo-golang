package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/horlerdipo/todo-golang/internal/auth"
	"github.com/horlerdipo/todo-golang/internal/todo"
	"github.com/horlerdipo/todo-golang/pkg"
	"gorm.io/gorm"
)

type Container struct {
	db            *gorm.DB
	AuthContainer *auth.Container
	TodoContainer *todo.Container
	EventBus      pkg.EventBus
}

func NewAppContainer(db *gorm.DB) *Container {
	eventBus := pkg.NewEventBus()
	return &Container{
		db:            db,
		AuthContainer: auth.NewContainer(db),
		TodoContainer: todo.NewContainer(db, eventBus),
		EventBus:      eventBus,
	}
}

func (container *Container) RegisterRoutes(r *chi.Mux) {
	container.AuthContainer.RegisterRoutes(r)
	container.TodoContainer.RegisterRoutes(r)
}

func (container *Container) RegisterListeners() {
	//container.AuthContainer.RegisterListeners(container.EventBus)
	container.TodoContainer.RegisterListeners(container.EventBus)
}
