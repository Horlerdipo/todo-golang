package todo

import (
	"errors"
	"github.com/horlerdipo/todo-golang/env"
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"github.com/horlerdipo/todo-golang/internal/enums"
	"golang.org/x/net/context"
	"log"
	"strconv"
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

func (service *Service) CreateTodo(ctx context.Context, createTodoDto *dtos.CreateTodoDTO) (uint, error) {
	if createTodoDto.Type == enums.Text {
		createTodoDto.Content = nil
	}

	todoId, err := service.TodoRepository.CreateTodo(ctx, createTodoDto)
	if err != nil {
		log.Println(err)
		return 0, errors.New("unable to create todo, please try again")
	}
	return todoId, nil
}

func (service *Service) DeleteTodo(ctx context.Context, todoId uint, userId uint) error {
	//check if user and to-do exists
	_, err := service.TodoRepository.FindTodoByUserId(ctx, todoId, userId)
	if err != nil {
		log.Println(err)
		return errors.New("todo does not exist")
	}

	err = service.TodoRepository.DeleteTodo(ctx, todoId)
	if err != nil {
		log.Println(err)
		return errors.New("unable to delete todo, please try again")
	}
	return nil
}

func (service *Service) AddChecklistItem(ctx context.Context, todoId uint, description string, userId uint) error {
	todo, err := service.TodoRepository.FindTodoByUserId(ctx, todoId, userId)
	if err != nil {
		log.Println(err)
		return errors.New("todo does not exist")
	}

	if todo.Type != enums.Checklist {
		return errors.New("only todos with type of checklists are supported")
	}

	_, err = service.TodoRepository.AddChecklistItem(ctx, todoId, description)
	return err
}

func (service *Service) DeleteChecklistItem(ctx context.Context, checklistId uint, todoId uint, userId uint) error {
	todo, err := service.TodoRepository.FindTodoByUserId(ctx, todoId, userId)
	if err != nil {
		log.Println(err)
		return errors.New("todo does not exist")
	}

	if todo.Type != enums.Checklist {
		return errors.New("only todos with type of checklists are supported")
	}

	return service.TodoRepository.DeleteChecklistItem(ctx, checklistId, todoId)
}

func (service *Service) UpdateChecklistItem(ctx context.Context, checklistId uint, description string, todoId uint, userId uint) (uint, error) {
	todo, err := service.TodoRepository.FindTodoByUserId(ctx, todoId, userId)
	if err != nil {
		log.Println(err)
		return 0, errors.New("todo does not exist")
	}

	if todo.Type != enums.Checklist {
		return 0, errors.New("only todos with type of checklists are supported")
	}

	return service.TodoRepository.UpdateChecklistItem(ctx, checklistId, todoId, description)
}

func (service *Service) UpdateChecklistItemStatus(ctx context.Context, checklistId uint, done bool, todoId uint, userId uint) (uint, error) {
	todo, err := service.TodoRepository.FindTodoByUserId(ctx, todoId, userId)
	if err != nil {
		log.Println(err)
		return 0, errors.New("todo does not exist")
	}

	log.Printf("Todo: +%v", todo)
	if todo.Type != enums.Checklist {
		return 0, errors.New("only todos with type of checklists are supported")
	}

	return service.TodoRepository.UpdateChecklistItemStatus(ctx, checklistId, todoId, done)
}

func (service *Service) PinTodo(ctx context.Context, todoId uint, userId uint) error {
	//check if the number of pinned is not more than 5
	maxPinnedTodos := env.FetchInt("MAXIMUM_PINNED_TODOS", 1)
	pinnedTodos := service.TodoRepository.CountPinnedTodos(ctx, userId)
	if int(pinnedTodos) >= maxPinnedTodos {
		return errors.New("you can only pin " + strconv.Itoa(maxPinnedTodos) + " todos")
	}

	//check if to-do belongs to user
	_, err := service.TodoRepository.FindTodoByUserId(ctx, todoId, userId)
	if err != nil {
		log.Println(err)
		return errors.New("todo does not exist")
	}

	return service.TodoRepository.PinTodo(ctx, todoId)
}

func (service *Service) UnPinTodo(ctx context.Context, todoId uint, userId uint) error {

	//check if to-do belongs to user
	_, err := service.TodoRepository.FindTodoByUserId(ctx, todoId, userId)
	if err != nil {
		log.Println(err)
		return errors.New("todo does not exist")
	}

	return service.TodoRepository.UnPinTodo(ctx, todoId)
}
