package dtos

import "github.com/horlerdipo/todo-golang/internal/enums"

type CreateTodoDTO struct {
	Title   string         `json:"title" validate:"required"`
	Content *string        `json:"content" validate:"required_if=Type text"`
	Type    enums.TodoType `json:"type" validate:"required,oneof=task text"`
	UserID  uint           `json:"user_id" validate:"-"`
}
