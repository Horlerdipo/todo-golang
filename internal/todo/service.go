package todo

import (
	"errors"
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"github.com/horlerdipo/todo-golang/internal/enums"
)

type Service struct {
	TodoRepository           database.TodoRepository
	TokenBlacklistRepository database.TokenBlacklistRepository
}

func NewService(todoRepository database.TodoRepository, blacklistRepository database.TokenBlacklistRepository) *Service {
	return &Service{
		todoRepository,
		blacklistRepository,
	}
}

func (service *Service) CreateTodo(createTodoDto *dtos.CreateTodoDTO) (uint, error) {
	if createTodoDto.Type == enums.Task {
		createTodoDto.Content = nil
	}

	todoId, err := service.TodoRepository.CreateTodo(createTodoDto)
	if err != nil {
		return 0, errors.New("unable to create todo, please try again")
	}
	return todoId, nil
}
