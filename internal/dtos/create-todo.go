package dtos

import "github.com/horlerdipo/todo-golang/internal/enums"

type CreateTodoDTO struct {
	Title     string         `json:"title" validate:"required"`
	Content   *string        `json:"content" validate:"required_if=Type text"`
	Type      enums.TodoType `json:"type" validate:"required,oneof=checklist text"`
	UserID    uint           `json:"user_id" validate:"-"`
	Checklist []string       `json:"checklist" validate:"required_if=Type checklist,omitempty,gt=0,dive,required"`
}

type ChecklistItem struct {
	Item string `json:"item" validate:"required"`
}

type ChecklistStatus struct {
	Done bool `json:"done" validate:"boolean"`
}
