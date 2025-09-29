package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/horlerdipo/todo-golang/internal/sse"
	"gorm.io/gorm"
)

type Container struct {
	AuthHandler *Handler
	AuthService *Service
	SSEService  *sse.Service
}

func NewContainer(db *gorm.DB, sseService *sse.Service) *Container {
	authService := NewService(
		database.NewUserRepository(db),
		database.NewTokenBlacklistRepository(db),
		sseService,
	)

	return &Container{
		AuthHandler: NewAuthHandler(authService),
		AuthService: authService,
		SSEService:  sseService,
	}
}

func (uc *Container) RegisterRoutes(r chi.Router) {
	uc.AuthHandler.RegisterRoutes(r)
}
