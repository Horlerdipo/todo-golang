package todo

import (
	"github.com/go-chi/chi/v5"
	"github.com/horlerdipo/todo-golang/internal/database"
	"gorm.io/gorm"
)

type Container struct {
	TodoService *Service
	TodoHandler *Handler
}

func NewContainer(db *gorm.DB) *Container {
	todoService := NewService(
		database.NewTodoRepository(db),
		database.NewTokenBlacklistRepository(db),
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
