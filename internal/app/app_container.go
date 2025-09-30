package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/horlerdipo/todo-golang/internal/auth"
	"github.com/horlerdipo/todo-golang/internal/sse"
	"github.com/horlerdipo/todo-golang/internal/todo"
	"github.com/horlerdipo/todo-golang/pkg"
	"gorm.io/gorm"
	"net/http"
)

type Container struct {
	db            *gorm.DB
	AuthContainer *auth.Container
	TodoContainer *todo.Container
	EventBus      pkg.EventBus
	SSEContainer  *sse.Container
}

func NewAppContainer(db *gorm.DB) *Container {
	eventBus := pkg.NewEventBus()
	sseContainer := sse.NewContainer(db)
	return &Container{
		db:            db,
		AuthContainer: auth.NewContainer(db, sseContainer.SSEService),
		TodoContainer: todo.NewContainer(db, eventBus, sseContainer.SSEService),
		EventBus:      eventBus,
		SSEContainer:  sseContainer,
	}
}

func (container *Container) RegisterRoutes(r *chi.Mux) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})
	container.AuthContainer.RegisterRoutes(r)
	container.TodoContainer.RegisterRoutes(r)
	container.SSEContainer.RegisterRoutes(r)
}

func (container *Container) RegisterListeners() {
	//container.AuthContainer.RegisterListeners(container.EventBus)
	container.TodoContainer.RegisterListeners(container.EventBus)
}
