package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/horlerdipo/todo-golang/env"
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/horlerdipo/todo-golang/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestPinTodo(t *testing.T) {
	//ARRANGE:
	tests := []struct {
		name               string
		setupFunc          func(t *testing.T) PinTodoSetupResponse
		expectedStatusCode int
		expectedMsg        string
		extraAssertions    func(t *testing.T, setupFuncResponse PinTodoSetupResponse)
	}{
		{
			name:               "success",
			setupFunc:          setupPinTodoSuccessfullyTest,
			expectedStatusCode: http.StatusOK,
			expectedMsg:        "",
			extraAssertions:    assertTodoPinned,
		},
		{
			name:               "unauthorized",
			setupFunc:          setupPinTodoUnauthorizedTest,
			expectedStatusCode: http.StatusUnauthorized,
			expectedMsg:        "Unauthenticated",
			extraAssertions:    assertPinTodoUnauthorized,
		},
		{
			name:               "maximum pinned todo",
			setupFunc:          setupPinTodoMaximumTest,
			expectedStatusCode: http.StatusBadRequest,
			expectedMsg:        "you can only pin " + env.FetchString("MAXIMUM_PINNED_TODOS") + " todos",
			extraAssertions:    assertPinTodoMaximum,
		},
		{
			name:               "todo not found",
			setupFunc:          setupPinTodoNotFoundTest,
			expectedStatusCode: http.StatusBadRequest,
			expectedMsg:        "todo does not exist",
			extraAssertions:    assertPinTodoNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupResponse := tt.setupFunc(t)
			url := fmt.Sprintf("%s/todos/%d/pin", TestServerInstance.Server.URL, setupResponse.Todo.ID)

			req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(nil))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+setupResponse.AuthToken)

			//ACT
			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			//ASSERT
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

func SetupPinTest(t *testing.T) (*database.User, string, *database.Todo) {
	ClearAllTables(t, TestServerInstance.DB)
	user := SeedUser(t, struct{}{})
	todo := SeedTodo(t, struct{}{}, user.ID)
	authToken := GenerateTestJwtToken(t, user.ID)
	return user, authToken, todo
}

type PinTodoSetupResponse struct {
	User      *database.User
	AuthToken string
	Todo      *database.Todo
}

func setupPinTodoSuccessfullyTest(t *testing.T) PinTodoSetupResponse {
	user, authToken, todo := SetupPinTest(t)
	return PinTodoSetupResponse{
		User:      user,
		AuthToken: authToken,
		Todo:      todo,
	}
}

func assertTodoPinned(t *testing.T, setupFuncResponse PinTodoSetupResponse) {
	newTodo := &database.Todo{}
	result := TestServerInstance.DB.Where("id = ?", setupFuncResponse.Todo.ID).First(&newTodo)
	assert.NoError(t, result.Error)
	assert.Equal(t, true, newTodo.Pinned)
}

func setupPinTodoUnauthorizedTest(t *testing.T) PinTodoSetupResponse {
	user, _, todo := SetupPinTest(t)
	return PinTodoSetupResponse{
		User:      user,
		AuthToken: "incorrect-token",
		Todo:      todo,
	}
}

func assertPinTodoUnauthorized(t *testing.T, setupFuncResponse PinTodoSetupResponse) {
	newTodo := &database.Todo{}
	result := TestServerInstance.DB.Where("id = ?", setupFuncResponse.Todo.ID).First(&newTodo)
	assert.NoError(t, result.Error)
	assert.Equal(t, false, newTodo.Pinned)
}

func setupPinTodoMaximumTest(t *testing.T) PinTodoSetupResponse {
	user, authToken, todo := SetupPinTest(t)

	maximumPinnedTodo := env.FetchInt("MAXIMUM_PINNED_TODOS", 1)
	for i := 0; i < maximumPinnedTodo; i++ {
		SeedTodo(t, database.Todo{
			Pinned: true,
		}, user.ID)
	}

	return PinTodoSetupResponse{
		User:      user,
		AuthToken: authToken,
		Todo:      todo,
	}
}

func assertPinTodoMaximum(t *testing.T, setupFuncResponse PinTodoSetupResponse) {
	newTodo := &database.Todo{}
	result := TestServerInstance.DB.Where("id = ?", setupFuncResponse.Todo.ID).First(&newTodo)
	assert.NoError(t, result.Error)
	assert.Equal(t, false, newTodo.Pinned)
}

func setupPinTodoNotFoundTest(t *testing.T) PinTodoSetupResponse {
	user, authToken, todo := SetupPinTest(t)
	todo.ID = todo.ID + 1
	return PinTodoSetupResponse{
		User:      user,
		AuthToken: authToken,
		Todo:      todo,
	}
}

func assertPinTodoNotFound(t *testing.T, setupFuncResponse PinTodoSetupResponse) {
	newTodo := &database.Todo{}
	result := TestServerInstance.DB.Where("id = ?", setupFuncResponse.Todo.ID-1).First(&newTodo)
	assert.NoError(t, result.Error)
	assert.Equal(t, false, newTodo.Pinned)
}
