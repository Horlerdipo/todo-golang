package events

type TodoCreatedEvent struct {
	TodoId uint
	UserId uint
}

func (event *TodoCreatedEvent) Name() string {
	return "todo.created"
}
