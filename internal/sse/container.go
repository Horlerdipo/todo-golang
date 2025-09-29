package sse

import (
	"github.com/go-chi/chi/v5"
	"github.com/horlerdipo/todo-golang/internal/database"
	"gorm.io/gorm"
)

type Container struct {
	SSEHandler *Handler
	SSEService *Service
}

func NewContainer(db *gorm.DB) *Container {
	tokenBlacklistRepository := database.NewTokenBlacklistRepository(db)
	service := NewService(tokenBlacklistRepository)

	return &Container{
		SSEHandler: NewHandler(service),
		SSEService: service,
	}
}

func (c *Container) RegisterRoutes(r chi.Router) {
	c.SSEHandler.RegisterRoutes(r)
}
