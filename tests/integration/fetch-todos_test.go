package integration

import (
	"github.com/horlerdipo/todo-golang/internal/database"
	"testing"
)

func setupFetchTodoTest(t *testing.T, todoCount int) (*database.User, string, []*database.Todo) {
	ClearAllTables(t, TestServerInstance.DB)
	user := SeedUser(t, struct{}{})
	authToken := GenerateTestJwtToken(t, user.ID)
	todos := make([]*database.Todo, todoCount)
	for i := 0; i < todoCount; i++ {
		SeedTodo(t, struct{}{}, user.ID)
	}
	return user, authToken, todos
}

func TestFetchTodos(t *testing.T) {}
