package events

type TodoCreatedEvent struct {
	TodoId uint
}

func (event *TodoCreatedEvent) Name() string {
	return "todo.created"
}
