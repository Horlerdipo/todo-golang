package database

import (
	"errors"
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"github.com/horlerdipo/todo-golang/internal/enums"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"log"
)

type TodoRepository interface {
	CreateTodo(ctx context.Context, createTodoDto *dtos.CreateTodoDTO) (uint, error)
	DeleteTodo(ctx context.Context, todoId uint) error
	FindTodoByUserId(ctx context.Context, todoId uint, userId uint) (*Todo, error)
	UpdateTodo(ctx context.Context, todoId uint, updateTodoDto *dtos.UpdateTodoDTO, deleteChecklist bool) error
	PinTodo(ctx context.Context, todoId uint) error
	UnPinTodo(ctx context.Context, todoId uint) error
	CountPinnedTodos(ctx context.Context, userId uint) int64

	AddChecklistItem(ctx context.Context, todoId uint, description string) (uint, error)
	DeleteChecklistItem(ctx context.Context, checklistId uint, todoId uint) error
	UpdateChecklistItem(ctx context.Context, checklistId uint, todoId uint, description string) (uint, error)
	UpdateChecklistItemStatus(ctx context.Context, checklistId uint, todoId uint, done bool) (uint, error)
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return todoRepository{db: db}
}

func (repo todoRepository) CreateTodo(ctx context.Context, createTodoDto *dtos.CreateTodoDTO) (uint, error) {
	todoModel := Todo{
		Content: createTodoDto.Content,
		Title:   createTodoDto.Title,
		Type:    createTodoDto.Type,
		UserID:  createTodoDto.UserID,
	}

	err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

func (repo todoRepository) DeleteTodo(ctx context.Context, todoId uint) error {
	err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx.Where("todo_id = ?", todoId).Delete(&Checklist{})
		tx.Delete(&Todo{}, todoId)
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo todoRepository) FindTodoByUserId(ctx context.Context, todoId uint, userId uint) (*Todo, error) {
	todo := Todo{}
	result := repo.db.WithContext(ctx).Where("user_id = ?", userId).Where("id = ?", todoId).First(&todo)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("todo not found")
		}
		return nil, result.Error
	}
	return &todo, nil
}

func (repo todoRepository) UpdateTodo(ctx context.Context, todoId uint, updateTodoDto *dtos.UpdateTodoDTO, deleteChecklist bool) error {

	var todo Todo
	result := repo.db.WithContext(ctx).Where("id = ?", todoId).First(&todo)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("todo not found")
		}
		log.Println("UpdateTodo error:", result.Error)
		return errors.New("unable to fetch todo while updating todo")
	}

	err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&todo).Select("Content", "Title", "Type").Where("id = ?", todo.ID).Updates(Todo{
			Content: updateTodoDto.Content,
			Title:   updateTodoDto.Title,
			Type:    updateTodoDto.Type,
		})

		if result.Error != nil {
			return result.Error
		}

		if deleteChecklist {
			tx.Model(&Checklist{}).Where("todo_id = ?", todoId).Delete(&Checklist{})
		}
		return nil
	})

	if err != nil {
		log.Printf("unable to update todo: %v", err)
		return errors.New("unable to update todo")
	}
	return nil
}

func (repo todoRepository) PinTodo(ctx context.Context, todoId uint) error {
	result := repo.db.WithContext(ctx).Model(&Todo{}).Where("id = ?", todoId).Update("pinned", true)
	if result.Error != nil {
		log.Println(result.Error)
		return errors.New("unable to pin todo, please try again")
	}
	return nil
}

func (repo todoRepository) UnPinTodo(ctx context.Context, todoId uint) error {
	result := repo.db.WithContext(ctx).Model(&Todo{}).Where("id = ?", todoId).Update("pinned", false)
	if result.Error != nil {
		log.Println(result.Error)
		return errors.New("unable to pin todo, please try again")
	}
	return nil
}

func (repo todoRepository) CountPinnedTodos(ctx context.Context, userId uint) int64 {
	var count int64
	repo.db.WithContext(ctx).Model(&Todo{}).Where("user_id = ?", userId).Count(&count)
	return count
}

func (repo todoRepository) AddChecklistItem(ctx context.Context, todoId uint, description string) (uint, error) {
	checklist := Checklist{
		Description: description,
		Done:        false,
		TodoID:      todoId,
	}

	result := repo.db.WithContext(ctx).Create(&checklist)
	if result.Error != nil {
		return 0, result.Error
	}
	return checklist.ID, nil
}

func (repo todoRepository) DeleteChecklistItem(ctx context.Context, checklistId uint, todoId uint) error {
	result := repo.db.WithContext(ctx).Where("id = ?", checklistId).Where("todo_id = ?", todoId).Delete(&Checklist{})
	if result.Error != nil {
		log.Print("Error while deleting checklist", result.Error)
		return errors.New("error while deleting checklist")
	}
	return nil
}

func (repo todoRepository) UpdateChecklistItem(ctx context.Context, checklistId uint, todoId uint, description string) (uint, error) {
	result := repo.db.WithContext(ctx).Model(&Checklist{}).Where("todo_id = ?", todoId).Where("id = ?", checklistId).Update("description", description)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, errors.New("checklist not found")
		}

		log.Println("Error while getting checklist", result.Error)
		return 0, errors.New("error while updating checklist")
	}
	return checklistId, nil
}

func (repo todoRepository) UpdateChecklistItemStatus(ctx context.Context, checklistId uint, todoId uint, done bool) (uint, error) {
	result := repo.db.WithContext(ctx).Model(&Checklist{}).Where("todo_id = ?", todoId).Where("id = ?", checklistId).Update("Done", done)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, errors.New("checklist not found")
		}

		log.Println("Error while getting checklist status", result.Error)
		return 0, errors.New("error while updating checklist status")
	}
	return checklistId, nil
}
