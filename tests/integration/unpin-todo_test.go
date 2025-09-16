package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/horlerdipo/todo-golang/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestUnpinTodo(t *testing.T) {
	//ARRANGE:
	tests := []struct {
		name               string
		setupFunc          func(t *testing.T) UnpinTodoSetupResponse
		expectedStatusCode int
		expectedMsg        string
		extraAssertions    func(t *testing.T, setupResponse UnpinTodoSetupResponse)
	}{
		{
			name:               "success",
			setupFunc:          setupUnPinTodoSuccess,
			expectedStatusCode: http.StatusOK,
			expectedMsg:        "",
			extraAssertions:    assertUnPinTodoSuccess,
		},
		{
			name:               "not found",
			setupFunc:          setupUnPinTodoNotFound,
			expectedStatusCode: http.StatusBadRequest,
			expectedMsg:        "todo does not exist",
			extraAssertions:    nil,
		},
		{
			name:               "unauthorized",
			setupFunc:          setupUnPinTodoUnauthorized,
			expectedStatusCode: http.StatusUnauthorized,
			expectedMsg:        "Unauthenticated",
			extraAssertions:    nil,
		},
		{
			name:               "different user todo",
			setupFunc:          setupUnPinTodoDifferentUser,
			expectedStatusCode: http.StatusBadRequest,
			expectedMsg:        "todo does not exist",
			extraAssertions:    assertUnPinTodoDifferentUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupResponse := tt.setupFunc(t)
			url := fmt.Sprintf("%s/todos/%d/unpin", TestServerInstance.Server.URL, setupResponse.Todo.ID)

			req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(nil))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+setupResponse.AuthToken)

			//ACT
			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			//ASSERT:
			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			if tt.expectedMsg != "" {
				var response utils.JsonResponse[interface{}]
				err := json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedMsg, response.Message)
			}

			if tt.extraAssertions != nil {
				tt.extraAssertions(t, setupResponse)
			}
		})
	}
}

type UnpinTodoSetupResponse struct {
	User      *database.User
	AuthToken string
	Todo      *database.Todo
}

func SetupUnPinTest(t *testing.T) (*database.User, string, *database.Todo) {
	ClearAllTables(t, TestServerInstance.DB)
	user := SeedUser(t, struct{}{})
	todo := SeedTodo(t, database.Todo{Pinned: true}, user.ID)
	authToken := GenerateTestJwtToken(t, user.ID)
	return user, authToken, todo
}

func setupUnPinTodoSuccess(t *testing.T) UnpinTodoSetupResponse {
	user, authToken, todo := SetupUnPinTest(t)
	return UnpinTodoSetupResponse{
		User:      user,
		AuthToken: authToken,
		Todo:      todo,
	}
}

func assertUnPinTodoSuccess(t *testing.T, setupResponse UnpinTodoSetupResponse) {
	newTodo := &database.Todo{}
	result := TestServerInstance.DB.Where("id = ?", setupResponse.Todo.ID).First(&newTodo)
	assert.NoError(t, result.Error)
	assert.Equal(t, false, newTodo.Pinned)
}

func setupUnPinTodoUnauthorized(t *testing.T) UnpinTodoSetupResponse {
	user, _, todo := SetupUnPinTest(t)
	return UnpinTodoSetupResponse{
		User:      user,
		AuthToken: "incorrect-auth-token",
		Todo:      todo,
	}
}

func setupUnPinTodoNotFound(t *testing.T) UnpinTodoSetupResponse {
	user, authToken, todo := SetupUnPinTest(t)
	todo.ID = todo.ID + 1
	return UnpinTodoSetupResponse{
		User:      user,
		AuthToken: authToken,
		Todo:      todo,
	}
}

func setupUnPinTodoDifferentUser(t *testing.T) UnpinTodoSetupResponse {
	user, _, todo := SetupUnPinTest(t)
	anotherUser := &database.User{}

	anotherUser = SeedUser(t, struct {
		Email string
	}{
		Email: "user2@gmail.com",
	})
	authToken := GenerateTestJwtToken(t, anotherUser.ID)

	return UnpinTodoSetupResponse{
		User:      user,
		AuthToken: authToken,
		Todo:      todo,
	}
}

func assertUnPinTodoDifferentUser(t *testing.T, setupResponse UnpinTodoSetupResponse) {
	newTodo := &database.Todo{}
	result := TestServerInstance.DB.Where("id = ?", setupResponse.Todo.ID).First(&newTodo)
	assert.NoError(t, result.Error)
	assert.Equal(t, true, newTodo.Pinned)
}
