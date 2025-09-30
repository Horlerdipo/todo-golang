package dtos

type SSEEventType string

const (
	TodoCreated      SSEEventType = "todoCreated"
	TodoDeleted      SSEEventType = "todoDeleted"
	TodoUpdated      SSEEventType = "todoUpdated"
	ChecklistAdded   SSEEventType = "checklistAdded"
	ChecklistDeleted SSEEventType = "checklistDeleted"
	ChecklistUpdated SSEEventType = "checklistUpdated"
)

type SSEData struct {
	Event SSEEventType `json:"event"`
	Data  interface{}  `json:"data"`
}
