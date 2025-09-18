package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/horlerdipo/todo-golang/utils"
	"net/http"
	"testing"
)

func TestUnpinTodo(t *testing.T) {
	tests := []struct {
		name               string
		setupFunc          func(t *testing.T) unpinTodoSetup
		expectedStatusCode int
		expectedMsg        string
		extraAssertions    func(t *testing.T, setup unpinTodoSetup)
	}{
		{
			name:               "success",
			setupFunc:          setupUnpinTodoSuccess,
			expectedStatusCode: http.StatusOK,
			expectedMsg:        "",
			extraAssertions:    assertUnpinTodoSuccess,
		},
		{
			name:               "not found",
			setupFunc:          setupUnpinTodoNotFound,
			expectedStatusCode: http.StatusBadRequest,
			expectedMsg:        "todo does not exist",
		},
		{
			name:               "unauthorized",
			setupFunc:          setupUnpinTodoUnauthorized,
			expectedStatusCode: http.StatusUnauthorized,
			expectedMsg:        "Unauthenticated",
		},
		{
			name:               "different user todo",
			setupFunc:          setupUnpinTodoDifferentUser,
			expectedStatusCode: http.StatusBadRequest,
			expectedMsg:        "todo does not exist",
			extraAssertions:    assertUnpinTodoDifferentUser,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			setup := tc.setupFunc(t)
			url := fmt.Sprintf("%s/todos/%d/unpin", TestServerInstance.Server.URL, setup.Todo.ID)

			req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(nil))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+setup.AuthToken)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("failed to perform request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tc.expectedStatusCode, resp.StatusCode)
			}

			if tc.expectedMsg != "" {
				var response utils.JsonResponse[interface{}]
				if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if response.Message != tc.expectedMsg {
					t.Errorf("expected message %q, got %q", tc.expectedMsg, response.Message)
				}
			}

			if tc.extraAssertions != nil {
				tc.extraAssertions(t, setup)
			}
		})
	}
}

type unpinTodoSetup struct {
	User      *database.User
	AuthToken string
	Todo      *database.Todo
}

func setupUnpinTest(t *testing.T) (*database.User, string, *database.Todo) {
	t.Helper()
	ClearAllTables(t, TestServerInstance.DB)
	user := SeedUser(t, struct{}{})
	todo := SeedTodo(t, database.Todo{Pinned: true}, user.ID)
	authToken := GenerateTestJwtToken(t, user.ID)
	return user, authToken, todo
}

func setupUnpinTodoSuccess(t *testing.T) unpinTodoSetup {
	t.Helper()
	user, authToken, todo := setupUnpinTest(t)
	return unpinTodoSetup{
		User:      user,
		AuthToken: authToken,
		Todo:      todo,
	}
}

func setupUnpinTodoUnauthorized(t *testing.T) unpinTodoSetup {
	t.Helper()
	user, _, todo := setupUnpinTest(t)
	return unpinTodoSetup{
		User:      user,
		AuthToken: "incorrect-auth-token",
		Todo:      todo,
	}
}

func setupUnpinTodoNotFound(t *testing.T) unpinTodoSetup {
	t.Helper()
	user, authToken, todo := setupUnpinTest(t)
	todo.ID = todo.ID + 1
	return unpinTodoSetup{
		User:      user,
		AuthToken: authToken,
		Todo:      todo,
	}
}

func setupUnpinTodoDifferentUser(t *testing.T) unpinTodoSetup {
	t.Helper()
	user, _, todo := setupUnpinTest(t)

	anotherUser := SeedUser(t, struct {
		Email string
	}{
		Email: "user2@gmail.com",
	})
	authToken := GenerateTestJwtToken(t, anotherUser.ID)

	return unpinTodoSetup{
		User:      user,
		AuthToken: authToken,
		Todo:      todo,
	}
}

func assertUnpinTodoSuccess(t *testing.T, setup unpinTodoSetup) {
	t.Helper()
	var todo database.Todo
	result := TestServerInstance.DB.Where("id = ?", setup.Todo.ID).First(&todo)
	if result.Error != nil {
		t.Fatalf("failed to fetch todo: %v", result.Error)
	}
	if todo.Pinned {
		t.Error("expected todo to be unpinned, but it was still pinned")
	}
}

func assertUnpinTodoDifferentUser(t *testing.T, setup unpinTodoSetup) {
	t.Helper()
	var todo database.Todo
	result := TestServerInstance.DB.Where("id = ?", setup.Todo.ID).First(&todo)
	if result.Error != nil {
		t.Fatalf("failed to fetch todo: %v", result.Error)
	}
	if !todo.Pinned {
		t.Error("expected todo to remain pinned, but it was unpinned")
	}
}
