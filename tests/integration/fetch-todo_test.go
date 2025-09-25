package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/horlerdipo/todo-golang/utils"
	"io"
	"net/http"
	"testing"
)

func SetupFetchTodoTest(t *testing.T) (*database.User, string, *database.Todo) {
	ClearAllTables(t, TestServerInstance.DB)
	user := SeedUser(t, struct{}{})
	todo := SeedTodo(t, struct{}{}, user.ID)
	authToken := GenerateTestJwtToken(t, user.ID)
	return user, authToken, todo
}

type FetchTodoSetupResponse struct {
	AuthToken string
	User      *database.User
	Todo      *database.Todo
}

func TestFetchTodo(t *testing.T) {
	tests := []struct {
		description        string
		setupFunc          func(t *testing.T) FetchTodoSetupResponse
		expectedStatusCode int
		expectedMsg        string
		extraAssertions    func(t *testing.T, setupFuncResponse FetchTodoSetupResponse, responseData utils.JsonResponse[database.Todo])
	}{
		{
			description:        "single todo can be fetched successfully",
			setupFunc:          fetchTodoSuccessfullySetup,
			expectedStatusCode: http.StatusOK,
			expectedMsg:        "Todo fetched successfully",
			extraAssertions:    fetchTodoSuccessfullyExtraAssertions,
		},
		{
			description:        "single todo and it's checklist can be fetched successfully",
			setupFunc:          fetchTodoAndChecklistSetup,
			expectedStatusCode: http.StatusOK,
			expectedMsg:        "Todo fetched successfully",
			extraAssertions:    fetchTodoAndChecklistExtraAssertions,
		},
		{
			description:        "returns 400 on unknown todo ID",
			setupFunc:          fetchTodoUnknownTodoIDSetup,
			expectedStatusCode: http.StatusBadRequest,
			expectedMsg:        "todo not found",
			extraAssertions:    nil,
		},
		{
			description:        "returns 400 on another user's todo ID",
			setupFunc:          fetchTodoUnauthorizedTodoSetup,
			expectedStatusCode: http.StatusBadRequest,
			expectedMsg:        "todo not found",
			extraAssertions:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(subTest *testing.T) {

			setup := tt.setupFunc(subTest)

			url := fmt.Sprintf("%s/todos/%d", TestServerInstance.Server.URL, setup.Todo.ID)
			req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(nil))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+setup.AuthToken)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatusCode {
				response, _ := io.ReadAll(resp.Body)
				t.Log(string(response))
				t.Errorf("expected status code %v, but got %v", tt.expectedStatusCode, resp.StatusCode)
			}

			if tt.expectedMsg != "" {
				var jsonResponse utils.JsonResponse[database.Todo]
				err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
				if err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if jsonResponse.Message != tt.expectedMsg {
					t.Errorf("expected message %s, but got %s", tt.expectedMsg, jsonResponse.Message)
				}

				if tt.extraAssertions != nil {
					tt.extraAssertions(subTest, setup, jsonResponse)
				}
			}
		})
	}
}

func fetchTodoSuccessfullyExtraAssertions(t *testing.T, setupResponse FetchTodoSetupResponse, responseData utils.JsonResponse[database.Todo]) {
	if setupResponse.Todo.ID != responseData.Data.ID {
		t.Errorf("expected fetched todo ID to be %v, is %v", setupResponse.Todo.ID, responseData.Data.ID)
	}

	if len(responseData.Data.Checklists) > 0 {
		t.Errorf("expected Checklists to be empty")
	}
}

func fetchTodoSuccessfullySetup(t *testing.T) FetchTodoSetupResponse {
	user, authToken, todo := SetupFetchTodoTest(t)
	return FetchTodoSetupResponse{
		AuthToken: authToken,
		User:      user,
		Todo:      todo,
	}
}

func fetchTodoAndChecklistExtraAssertions(t *testing.T, setupResponse FetchTodoSetupResponse, responseData utils.JsonResponse[database.Todo]) {
	if setupResponse.Todo.ID != responseData.Data.ID {
		t.Errorf("expected fetched todo ID to be %v, is %v", setupResponse.Todo.ID, responseData.Data.ID)
	}

	if len(responseData.Data.Checklists) != 2 {
		t.Errorf("expected Checklists to be 2")
	}
}

func fetchTodoAndChecklistSetup(t *testing.T) FetchTodoSetupResponse {
	user, authToken, todo := SetupFetchTodoTest(t)
	SeedChecklist(t, struct{}{}, todo.ID)
	SeedChecklist(t, struct{}{}, todo.ID)
	return FetchTodoSetupResponse{
		AuthToken: authToken,
		User:      user,
		Todo:      todo,
	}
}

func fetchTodoUnknownTodoIDSetup(t *testing.T) FetchTodoSetupResponse {
	user, authToken, todo := SetupFetchTodoTest(t)
	todo.ID = todo.ID - 1
	return FetchTodoSetupResponse{
		AuthToken: authToken,
		User:      user,
		Todo:      todo,
	}
}

func fetchTodoUnauthorizedTodoSetup(t *testing.T) FetchTodoSetupResponse {
	user, authToken, _ := SetupFetchTodoTest(t)
	anotherUser := SeedUser(t, struct{}{})
	anotherTodo := SeedTodo(t, struct{}{}, anotherUser.ID)

	return FetchTodoSetupResponse{
		AuthToken: authToken,
		User:      user,
		Todo:      anotherTodo,
	}
}
