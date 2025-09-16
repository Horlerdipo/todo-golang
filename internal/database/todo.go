package database

import (
	"github.com/horlerdipo/todo-golang/internal/enums"
)

type Todo struct {
	Model
	Title   string         `json:"title"`
	Content *string        `json:"content"`
	Type    enums.TodoType `json:"type"`
	UserID  uint           `json:"user_id"`
	User    User           `gorm:"constraint:OnDelete:CASCADE" json:"user"`
	Pinned  bool           `gorm:"default:false" json:"pinned"`
}
