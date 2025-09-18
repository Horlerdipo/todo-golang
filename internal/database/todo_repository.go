package database

import (
	"errors"
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"github.com/horlerdipo/todo-golang/internal/enums"
	"gorm.io/gorm"
	"log"
)

type TodoRepository interface {
	CreateTodo(createTodoDto *dtos.CreateTodoDTO) (uint, error)
	DeleteTodo(todoId uint) error
	FindTodoByUserId(todoId uint, userId uint) (*Todo, error)
	PinTodo(todoId uint) error
	UnPinTodo(todoId uint) error
	CountPinnedTodos(userId uint) int64

	AddChecklistItem(todoId uint, description string) (uint, error)
	DeleteChecklistItem(checklistId uint, todoId uint) error
	UpdateChecklistItem(checklistId uint, todoId uint, description string) (uint, error)
	UpdateChecklistItemStatus(checklistId uint, todoId uint, done bool) (uint, error)
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

	err := repo.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&todoModel)
		if result.Error != nil {
			return result.Error
		}

		if createTodoDto.Type == enums.Checklist {
			var checklists []Checklist
			for _, checklist := range createTodoDto.Checklist {
				checklistModel := Checklist{
					Description: checklist,
					Done:        false,
					TodoID:      todoModel.ID,
				}
				checklists = append(checklists, checklistModel)
			}
			result = tx.Create(&checklists)
			if result.Error != nil {
				return result.Error
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}
	return todoModel.ID, nil
}

func (repo todoRepository) DeleteTodo(todoId uint) error {
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		tx.Where("todo_id = ?", todoId).Delete(&Checklist{})
		tx.Delete(&Todo{}, todoId)
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo todoRepository) FindTodoByUserId(todoId uint, userId uint) (*Todo, error) {
	todo := Todo{}
	result := repo.db.Where("user_id = ?", userId).Where("id = ?", todoId).First(&todo)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("todo not found")
		}

		return nil, result.Error
	}
	return &todo, nil
}

func (repo todoRepository) PinTodo(todoId uint) error {
	result := repo.db.Model(&Todo{}).Where("id = ?", todoId).Update("pinned", true)
	if result.Error != nil {
		log.Println(result.Error)
		return errors.New("unable to pin todo, please try again")
	}
	return nil
}

func (repo todoRepository) UnPinTodo(todoId uint) error {
	result := repo.db.Model(&Todo{}).Where("id = ?", todoId).Update("pinned", false)
	if result.Error != nil {
		log.Println(result.Error)
		return errors.New("unable to pin todo, please try again")
	}
	return nil
}

func (repo todoRepository) CountPinnedTodos(userId uint) int64 {
	var count int64
	repo.db.Model(&Todo{}).Where("user_id = ?", userId).Count(&count)
	return count
}

func (repo todoRepository) AddChecklistItem(todoId uint, description string) (uint, error) {
	checklist := Checklist{
		Description: description,
		Done:        false,
		TodoID:      todoId,
	}

	result := repo.db.Create(&checklist)
	if result.Error != nil {
		return 0, result.Error
	}
	return checklist.ID, nil
}

func (repo todoRepository) DeleteChecklistItem(checklistId uint, todoId uint) error {
	result := repo.db.Where("id = ?", checklistId).Where("todo_id = ?", todoId).Delete(&Checklist{})
	if result.Error != nil {
		log.Print("Error while deleting checklist", result.Error)
		return errors.New("error while deleting checklist")
	}
	return nil
}

func (repo todoRepository) UpdateChecklistItem(checklistId uint, todoId uint, description string) (uint, error) {
	result := repo.db.Model(&Checklist{}).Where("todo_id = ?", todoId).Where("id = ?", checklistId).Update("description", description)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, errors.New("checklist not found")
		}

		log.Println("Error while getting checklist", result.Error)
		return 0, errors.New("error while updating checklist")
	}
	return checklistId, nil
}

func (repo todoRepository) UpdateChecklistItemStatus(checklistId uint, todoId uint, done bool) (uint, error) {
	result := repo.db.Model(&Checklist{}).Where("todo_id = ?", todoId).Where("id = ?", checklistId).Update("Done", done)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, errors.New("checklist not found")
		}

		log.Println("Error while getting checklist status", result.Error)
		return 0, errors.New("error while updating checklist status")
	}
	return checklistId, nil
}
