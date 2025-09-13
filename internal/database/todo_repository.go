package database

import (
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"gorm.io/gorm"
)

type TodoRepository interface {
	CreateTodo(createTodoDto *dtos.CreateTodoDTO) (uint, error)
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return todoRepository{db: db}
}

func (repo todoRepository) CreateTodo(createTodoDto *dtos.CreateTodoDTO) (uint, error) {
	todoModel := Todo{
		Content: createTodoDto.Content,
		Title:   createTodoDto.Title,
		Type:    createTodoDto.Type,
		UserID:  createTodoDto.UserID,
	}

	result := repo.db.Create(&todoModel)
	if result.Error != nil {
		return 0, result.Error
	}
	return todoModel.ID, nil
}
