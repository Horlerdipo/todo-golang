package todo

import (
	"github.com/go-chi/chi/v5"
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/horlerdipo/todo-golang/pkg"
	"gorm.io/gorm"
)

type Container struct {
	TodoService *Service
	TodoHandler *Handler
}

func NewContainer(db *gorm.DB, bus pkg.EventBus) *Container {
	todoService := NewService(
		database.NewTodoRepository(db),
		database.NewTokenBlacklistRepository(db),
		bus,
	)

	todoHandler := NewHandler(todoService)

	return &Container{
		TodoService: todoService,
		TodoHandler: todoHandler,
	}
}

func (uc *Container) RegisterRoutes(r chi.Router) {
	uc.TodoHandler.RegisterRoutes(r)
}

func (uc *Container) RegisterListeners(bus pkg.EventBus) {
	bus.Subscribe("todo.created", NewTodoCreatedListener(uc.TodoService.TodoRepository))
}
