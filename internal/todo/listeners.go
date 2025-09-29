package todo

import (
	"fmt"
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"github.com/horlerdipo/todo-golang/internal/events"
	"github.com/horlerdipo/todo-golang/internal/sse"
	"github.com/horlerdipo/todo-golang/pkg"
)

type TodoCreatedListener struct {
	TodoRepository database.TodoRepository
	SSEService     *sse.Service
}

func (listener *TodoCreatedListener) Handle(event pkg.Event) {
	e := event.(*events.TodoCreatedEvent)
	fmt.Printf("Todo created listener triggered by %v with user id %v", e.TodoId, e.UserId)
	message := dtos.SSEData{
		Event: dtos.TodoCreated,
		Data:  e.TodoId,
	}
	listener.SSEService.SendMessage(e.UserId, message)
}

func NewTodoCreatedListener(repository database.TodoRepository, sseService *sse.Service) *TodoCreatedListener {
	return &TodoCreatedListener{
		TodoRepository: repository,
		SSEService:     sseService,
	}
}
