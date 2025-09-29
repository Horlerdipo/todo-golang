package dtos

type SSEEventType string

const (
	TodoCreated SSEEventType = "todoCreated"
	TodoDeleted SSEEventType = "todoDeleted"
)

type SSEData struct {
	Event SSEEventType
	Data  interface{}
}
