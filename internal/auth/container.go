package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/horlerdipo/todo-golang/internal/users"
	"gorm.io/gorm"
)

type Container struct {
	AuthHandler *Handler
	AuthService *Service
}

func NewContainer(db *gorm.DB) *Container {
	authService := NewService(users.NewUserRepository(db))

	return &Container{
		AuthHandler: NewAuthHandler(authService),
		AuthService: authService,
	}
}

func (uc *Container) RegisterRoutes(r chi.Router) {
	uc.AuthHandler.RegisterRoutes(r)
}
