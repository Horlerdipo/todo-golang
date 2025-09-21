package database

type Checklist struct {
	Model
	Description string `json:"description"`
	Done        bool   `json:"done"`
	TodoID      uint   `json:"todo_id"`
	Todo        Todo   `gorm:"constraint:OnDelete:CASCADE" json:"-"`
}
