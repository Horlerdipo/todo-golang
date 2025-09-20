package todo

import (
	"fmt"
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/horlerdipo/todo-golang/internal/events"
	"github.com/horlerdipo/todo-golang/pkg"
)

type TodoCreatedListener struct {
	TodoRepository database.TodoRepository
}

func (listener *TodoCreatedListener) Handle(event pkg.Event) {
	e := event.(*events.TodoCreatedEvent)
	fmt.Printf("Todo created listener triggered by %v", e.TodoId)
}

func NewTodoCreatedListener(repository database.TodoRepository) *TodoCreatedListener {
	return &TodoCreatedListener{TodoRepository: repository}
}
