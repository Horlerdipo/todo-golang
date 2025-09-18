package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/horlerdipo/todo-golang/internal/enums"
	"github.com/horlerdipo/todo-golang/utils"
	"net/http"
	"testing"
)

func SetupAddChecklistTest(t *testing.T) (*database.User, string, *database.Todo) {
	ClearAllTables(t, TestServerInstance.DB)
	user := SeedUser(t, struct{}{})
	todo := SeedTodo(t, struct {
		Type enums.TodoType
	}{
		Type: enums.Checklist,
	}, user.ID)
	authToken := GenerateTestJwtToken(t, user.ID)
	return user, authToken, todo
}

type AddItemToTodoChecklistSetupResponse struct {
	AuthToken     string
	User          *database.User
	Todo          *database.Todo
	ChecklistItem string
}

func TestAddChecklistToTodo(t *testing.T) {
	tests := []struct {
		description        string
		setupFunc          func(t *testing.T) AddItemToTodoChecklistSetupResponse
		expectedStatusCode int
		expectedMsg        string
		extraAssertions    func(t *testing.T, setupFuncResponse AddItemToTodoChecklistSetupResponse)
	}{
		{
			description:        "checklist can be added to todo successfully",
			setupFunc:          addItemToTodoChecklistSuccessfullySetup,
			expectedStatusCode: http.StatusCreated,
			expectedMsg:        "",
			extraAssertions:    addItemToTodoChecklistSuccessfullyExtraAssertions,
		},
		{
			description:        "checklist can only be added to Todos with type of checklist",
			setupFunc:          preventChecklistAdditionSetup,
			expectedStatusCode: http.StatusBadRequest,
			expectedMsg:        "only todos with type of checklists are supported",
			extraAssertions:    preventChecklistAdditionExtraAssertions,
		},
		{
			description:        "checklist returns validation error",
			setupFunc:          addItemToTodoChecklistValidationErrorSetup,
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedMsg:        "",
			extraAssertions:    addItemToTodoChecklistValidationErrorExtraAssertions,
		},
		{
			description:        "checklist returns error on unknown Todo",
			setupFunc:          addItemToTodoChecklistUnknownTodoSetup,
			expectedStatusCode: http.StatusBadRequest,
			expectedMsg:        "",
			extraAssertions:    addItemToTodoChecklistUnknownTodoExtraAssertions,
		},
		{
			description:        "checklist returns error when TODO ID is another user's",
			setupFunc:          addItemToTodoChecklistAnotherUserTodoSetup,
			expectedStatusCode: http.StatusBadRequest,
			expectedMsg:        "",
			extraAssertions:    addItemToTodoChecklistAnotherUserTodoExtraAssertions,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(subTest *testing.T) {

			setup := tt.setupFunc(subTest)

			jsonRequest, err := json.Marshal(map[string]string{
				"item": setup.ChecklistItem,
			})
			if err != nil {
				t.Fatal("Failed to marshal request")
			}

			url := fmt.Sprintf("%s/todos/%d/checklist", TestServerInstance.Server.URL, setup.Todo.ID)
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonRequest))
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
				t.Errorf("expected status code %v, but got %v", tt.expectedStatusCode, resp.StatusCode)
			}
			if tt.expectedMsg != "" {
				var jsonResponse utils.JsonResponse[interface{}]
				err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
				if err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if jsonResponse.Message != tt.expectedMsg {
					t.Errorf("expected message %s, but got %s", tt.expectedMsg, jsonResponse.Message)
				}

				tt.extraAssertions(subTest, setup)
			}
		})
	}
}

func addItemToTodoChecklistSuccessfullySetup(t *testing.T) AddItemToTodoChecklistSetupResponse {
	t.Helper()
	user, authToken, todo := SetupAddChecklistTest(t)
	return AddItemToTodoChecklistSetupResponse{
		AuthToken:     authToken,
		User:          user,
		Todo:          todo,
		ChecklistItem: "Testing",
	}
}

func addItemToTodoChecklistSuccessfullyExtraAssertions(t *testing.T, setup AddItemToTodoChecklistSetupResponse) {
	t.Helper()

	checklist := database.Checklist{}
	result := TestServerInstance.DB.Where("").First(&checklist)
	if result.Error != nil {
		t.Errorf("failed to find checklist: %v", result.Error)
	}

	if checklist.TodoID != setup.Todo.ID {
		t.Errorf("expected todo id %v, but got %v", setup.Todo.ID, checklist.TodoID)
	}

	if checklist.Description != setup.ChecklistItem {
		t.Errorf("expected checklist description %v, but got %v", setup.ChecklistItem, checklist.Description)
	}

	if checklist.Done {
		t.Error("expected checklist not to be marked as done, but marked as done")
	}
}

func addItemToTodoChecklistValidationErrorSetup(t *testing.T) AddItemToTodoChecklistSetupResponse {
	t.Helper()
	user, authToken, todo := SetupAddChecklistTest(t)
	return AddItemToTodoChecklistSetupResponse{
		AuthToken:     authToken,
		User:          user,
		Todo:          todo,
		ChecklistItem: "",
	}
}

func addItemToTodoChecklistValidationErrorExtraAssertions(t *testing.T, setup AddItemToTodoChecklistSetupResponse) {
	t.Helper()
	checklist := database.Checklist{}
	result := TestServerInstance.DB.Where("").First(&checklist)
	if result.Error == nil {
		t.Errorf("failed to find checklist: %v", result.Error)
	}
}

func addItemToTodoChecklistUnknownTodoSetup(t *testing.T) AddItemToTodoChecklistSetupResponse {
	t.Helper()
	user, authToken, todo := SetupAddChecklistTest(t)
	todo.ID = todo.ID + 1
	return AddItemToTodoChecklistSetupResponse{
		AuthToken:     authToken,
		User:          user,
		Todo:          todo,
		ChecklistItem: "Testing testing",
	}
}

func addItemToTodoChecklistUnknownTodoExtraAssertions(t *testing.T, setup AddItemToTodoChecklistSetupResponse) {
	t.Helper()
	checklist := database.Checklist{}
	result := TestServerInstance.DB.Where("").First(&checklist)
	if result.Error == nil {
		t.Errorf("failed to find checklist: %v", result.Error)
	}
}

func addItemToTodoChecklistAnotherUserTodoSetup(t *testing.T) AddItemToTodoChecklistSetupResponse {
	t.Helper()
	user, authToken, _ := SetupAddChecklistTest(t)
	anotherUser := SeedUser(t, struct{}{})
	anotherTodo := SeedTodo(t, struct{}{}, anotherUser.ID)
	return AddItemToTodoChecklistSetupResponse{
		AuthToken:     authToken,
		User:          user,
		Todo:          anotherTodo,
		ChecklistItem: "Testing testing",
	}
}

func addItemToTodoChecklistAnotherUserTodoExtraAssertions(t *testing.T, setup AddItemToTodoChecklistSetupResponse) {
	t.Helper()
	checklist := database.Checklist{}
	result := TestServerInstance.DB.Where("").First(&checklist)
	if result.Error == nil {
		t.Errorf("failed to find checklist: %v", result.Error)
	}
}

func preventChecklistAdditionSetup(t *testing.T) AddItemToTodoChecklistSetupResponse {
	t.Helper()
	user, authToken, todo := SetupAddChecklistTest(t)
	TestServerInstance.DB.Model(&todo).Where("id", todo.ID).Update("type", enums.Text)
	return AddItemToTodoChecklistSetupResponse{
		AuthToken:     authToken,
		User:          user,
		Todo:          todo,
		ChecklistItem: "Testing testing",
	}
}

func preventChecklistAdditionExtraAssertions(t *testing.T, setup AddItemToTodoChecklistSetupResponse) {
	t.Helper()
	checklist := database.Checklist{}
	result := TestServerInstance.DB.Where("todo_id", setup.Todo.ID).First(&checklist)
	if result.Error == nil {
		t.Errorf("failed to find checklist: %v", result.Error)
	}
}
