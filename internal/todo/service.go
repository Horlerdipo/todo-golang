package todo

import (
	"errors"
	"github.com/horlerdipo/todo-golang/env"
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"github.com/horlerdipo/todo-golang/internal/enums"
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

func (service *Service) CreateTodo(createTodoDto *dtos.CreateTodoDTO) (uint, error) {
	if createTodoDto.Type == enums.Text {
		createTodoDto.Content = nil
	}

	todoId, err := service.TodoRepository.CreateTodo(createTodoDto)
	if err != nil {
		log.Println(err)
		return 0, errors.New("unable to create todo, please try again")
	}
	return todoId, nil
}

func (service *Service) DeleteTodo(todoId uint, userId uint) error {
	//check if user and to-do exists
	_, err := service.TodoRepository.FindTodoByUserId(todoId, userId)
	if err != nil {
		log.Println(err)
		return errors.New("todo does not exist")
	}

	err = service.TodoRepository.DeleteTodo(todoId)
	if err != nil {
		log.Println(err)
		return errors.New("unable to delete todo, please try again")
	}
	return nil
}

func (service *Service) AddChecklistItem(todoId uint, description string, userId uint) error {
	todo, err := service.TodoRepository.FindTodoByUserId(todoId, userId)
	if err != nil {
		log.Println(err)
		return errors.New("todo does not exist")
	}

	if todo.Type != enums.Checklist {
		return errors.New("only todos with type of checklists are supported")
	}

	_, err = service.TodoRepository.AddChecklistItem(todoId, description)
	return err
}

func (service *Service) DeleteChecklistItem(checklistId uint, todoId uint, userId uint) error {
	todo, err := service.TodoRepository.FindTodoByUserId(todoId, userId)
	if err != nil {
		log.Println(err)
		return errors.New("todo does not exist")
	}

	if todo.Type != enums.Checklist {
		return errors.New("only todos with type of checklists are supported")
	}

	return service.TodoRepository.DeleteChecklistItem(checklistId, todoId)
}

func (service *Service) UpdateChecklistItem(checklistId uint, description string, todoId uint, userId uint) (uint, error) {
	todo, err := service.TodoRepository.FindTodoByUserId(todoId, userId)
	if err != nil {
		log.Println(err)
		return 0, errors.New("todo does not exist")
	}

	if todo.Type != enums.Checklist {
		return 0, errors.New("only todos with type of checklists are supported")
	}

	return service.TodoRepository.UpdateChecklistItem(checklistId, todoId, description)
}

func (service *Service) UpdateChecklistItemStatus(checklistId uint, done bool, todoId uint, userId uint) (uint, error) {
	todo, err := service.TodoRepository.FindTodoByUserId(todoId, userId)
	if err != nil {
		log.Println(err)
		return 0, errors.New("todo does not exist")
	}

	log.Printf("Todo: +%v", todo)
	if todo.Type != enums.Checklist {
		return 0, errors.New("only todos with type of checklists are supported")
	}

	return service.TodoRepository.UpdateChecklistItemStatus(checklistId, todoId, done)
}

func (service *Service) PinTodo(todoId uint, userId uint) error {
	//check if the number of pinned is not more than 5
	maxPinnedTodos := env.FetchInt("MAXIMUM_PINNED_TODOS", 1)
	pinnedTodos := service.TodoRepository.CountPinnedTodos(userId)
	if int(pinnedTodos) >= maxPinnedTodos {
		return errors.New("you can only pin " + strconv.Itoa(maxPinnedTodos) + " todos")
	}

	//check if to-do belongs to user
	_, err := service.TodoRepository.FindTodoByUserId(todoId, userId)
	if err != nil {
		log.Println(err)
		return errors.New("todo does not exist")
	}

	return service.TodoRepository.PinTodo(todoId)
}

func (service *Service) UnPinTodo(todoId uint, userId uint) error {

	//check if to-do belongs to user
	_, err := service.TodoRepository.FindTodoByUserId(todoId, userId)
	if err != nil {
		log.Println(err)
		return errors.New("todo does not exist")
	}

	return service.TodoRepository.UnPinTodo(todoId)
}
