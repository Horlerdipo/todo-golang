package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"github.com/horlerdipo/todo-golang/internal/enums"
	"github.com/horlerdipo/todo-golang/utils"
	"io"
	"net/http"
	"testing"
)

var updateTodoContent string = "random random string"

var updateTodoRequest dtos.UpdateTodoDTO = dtos.UpdateTodoDTO{
	Title:   "Hello World",
	Content: &updateTodoContent,
	Type:    enums.Text,
}

func setupUpdateTodoTest(t *testing.T) (*database.User, string, *database.Todo) {
	t.Helper()
	ClearAllTables(t, TestServerInstance.DB)
	user := SeedUser(t, struct{}{})
	authToken := GenerateTestJwtToken(t, user.ID)
	todo := SeedTodo(t, struct{}{}, user.ID)
	return user, authToken, todo
}

type UpdateTodoSetupResponse struct {
	User       *database.User
	AuthToken  string
	Todo       *database.Todo
	RequestDto dtos.UpdateTodoDTO
}

func TestUpdateTodo(t *testing.T) {
	tests := []struct {
		description        string
		setupFunc          func(t *testing.T) UpdateTodoSetupResponse
		expectedStatusCode int
		expectedMsg        string
		extraAssertions    func(t *testing.T, setupFuncResponse UpdateTodoSetupResponse)
	}{
		{
			description:        "Should update todo successfully",
			setupFunc:          updateTodoSuccessfullySetup,
			expectedStatusCode: http.StatusNoContent,
			expectedMsg:        "",
			extraAssertions:    updateTodoSuccessfullyExtraAssertions,
		},
		{
			description:        "Should update todo successfully and change Todo Content to Nil",
			setupFunc:          changeTodoContentToNilSuccessfullySetup,
			expectedStatusCode: http.StatusNoContent,
			expectedMsg:        "",
			extraAssertions:    changeTodoContentToNilSuccessfullyExtraAssertions,
		},
		{
			description:        "Should update todo successfully and delete checklists associated with the todo",
			setupFunc:          clearTodoChecklistSuccessfullySetup,
			expectedStatusCode: http.StatusNoContent,
			expectedMsg:        "",
			extraAssertions:    clearTodoChecklistSuccessfullyExtraAssertions,
		},
		{
			description:        "Should return validation error",
			setupFunc:          updateTodoValidationErrorSetup,
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedMsg:        "",
			extraAssertions:    nil,
		},
		{
			description:        "Should return error on incorrect todo ID",
			setupFunc:          returnErrorOnIncorrectTodoIDSetup,
			expectedStatusCode: http.StatusBadRequest,
			expectedMsg:        "",
			extraAssertions:    returnErrorOnIncorrectTodoIDExtraAssertions,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			setup := test.setupFunc(t)
			jsonRequest, err := json.Marshal(setup.RequestDto)
			if err != nil {
				t.Fatal("Failed to marshal request")
			}

			url := fmt.Sprintf("%s/todos/%d", TestServerInstance.Server.URL, setup.Todo.ID)
			req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonRequest))
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

			if resp.StatusCode != test.expectedStatusCode {
				body, _ := io.ReadAll(resp.Body)
				t.Logf("request body for %v : %v", test.description, string(jsonRequest))
				t.Log(string(body))
				t.Errorf("expected status code %v, but got %v", test.expectedStatusCode, resp.StatusCode)
			}
			if test.expectedMsg != "" {
				var jsonResponse utils.JsonResponse[interface{}]
				err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
				if err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if jsonResponse.Message != test.expectedMsg {
					t.Errorf("expected message %s, but got %s", test.expectedMsg, jsonResponse.Message)
				}

			}

			if test.extraAssertions != nil {
				test.extraAssertions(t, setup)
			}

		})
	}
}

func updateTodoSuccessfullySetup(t *testing.T) UpdateTodoSetupResponse {
	t.Helper()
	user, authToken, todo := setupUpdateTodoTest(t)
	request := dtos.UpdateTodoDTO{
		Title:   "Hello World",
		Content: &updateTodoContent,
		Type:    enums.Text,
	}

	return UpdateTodoSetupResponse{
		User:       user,
		AuthToken:  authToken,
		Todo:       todo,
		RequestDto: request,
	}
}

func updateTodoSuccessfullyExtraAssertions(t *testing.T, setup UpdateTodoSetupResponse) {
	t.Helper()
	todo := database.Todo{}
	result := TestServerInstance.DB.Where("id = ?", setup.Todo.ID).First(&todo)
	if result.Error != nil {
		t.Fatalf("unable to query todo: %v", result.Error)
	}

	if setup.RequestDto.Title != todo.Title {
		t.Errorf("expected %v, got %v", setup.RequestDto.Title, todo.Title)
	}

	if setup.RequestDto.Type != todo.Type {
		t.Errorf("expected %v, got %v", setup.RequestDto.Type, todo.Type)
	}

	if *setup.RequestDto.Content != *todo.Content {
		t.Errorf("expected %v, got %v", setup.RequestDto.Content, todo.Content)
	}
}

func changeTodoContentToNilSuccessfullySetup(t *testing.T) UpdateTodoSetupResponse {
	t.Helper()
	user, authToken, todo := setupUpdateTodoTest(t)
	request := dtos.UpdateTodoDTO{
		Title:   "Hello World",
		Content: nil,
		Type:    enums.Checklist,
	}

	return UpdateTodoSetupResponse{
		User:       user,
		AuthToken:  authToken,
		Todo:       todo,
		RequestDto: request,
	}
}

func changeTodoContentToNilSuccessfullyExtraAssertions(t *testing.T, setup UpdateTodoSetupResponse) {
	t.Helper()
	todo := database.Todo{}
	result := TestServerInstance.DB.Where("id = ?", setup.Todo.ID).First(&todo)
	if result.Error != nil {
		t.Fatalf("unable to query todo: %v", result.Error)
	}

	if setup.RequestDto.Title != todo.Title {
		t.Errorf("expected %v, got %v", setup.RequestDto.Title, todo.Title)
	}

	if setup.RequestDto.Type != todo.Type {
		t.Errorf("expected %v, got %v", setup.RequestDto.Type, todo.Type)
	}

	if todo.Content != nil {
		t.Errorf("expected %v, got %v", nil, todo.Content)
	}
}

func clearTodoChecklistSuccessfullySetup(t *testing.T) UpdateTodoSetupResponse {
	t.Helper()
	user, authToken, todo := setupUpdateTodoTest(t)
	result := TestServerInstance.DB.Model(&todo).Where("id = ?", todo.ID).Update("type", enums.Checklist)
	if result.Error != nil {
		t.Fatalf("unable to update todo type: %v", result.Error)
	}

	SeedChecklist(t, struct{}{}, todo.ID)

	newContent := "new content here"
	request := dtos.UpdateTodoDTO{
		Title:   "Hello World",
		Content: &newContent,
		Type:    enums.Text,
	}

	return UpdateTodoSetupResponse{
		User:       user,
		AuthToken:  authToken,
		Todo:       todo,
		RequestDto: request,
	}
}

func clearTodoChecklistSuccessfullyExtraAssertions(t *testing.T, setup UpdateTodoSetupResponse) {
	t.Helper()
	todo := database.Todo{}
	result := TestServerInstance.DB.Where("id = ?", setup.Todo.ID).First(&todo)
	if result.Error != nil {
		t.Fatalf("unable to query todo: %v", result.Error)
	}

	if setup.RequestDto.Title != todo.Title {
		t.Errorf("expected %v, got %v", setup.RequestDto.Title, todo.Title)
	}

	if setup.RequestDto.Type != todo.Type {
		t.Errorf("expected %v, got %v", setup.RequestDto.Type, todo.Type)
	}

	if *todo.Content != *setup.RequestDto.Content {
		t.Errorf("expected %v, got %v", *setup.RequestDto.Content, *todo.Content)
	}

	checklist := database.Checklist{}
	result = TestServerInstance.DB.Where("todo_id = ?", setup.Todo.ID).First(&checklist)
	if result.Error == nil {
		t.Error("checklist should not exist in the database")
	}
}

func updateTodoValidationErrorSetup(t *testing.T) UpdateTodoSetupResponse {
	t.Helper()
	user, authToken, todo := setupUpdateTodoTest(t)
	request := dtos.UpdateTodoDTO{}

	return UpdateTodoSetupResponse{
		User:       user,
		AuthToken:  authToken,
		Todo:       todo,
		RequestDto: request,
	}
}

func returnErrorOnIncorrectTodoIDSetup(t *testing.T) UpdateTodoSetupResponse {
	t.Helper()
	user, authToken, todo := setupUpdateTodoTest(t)
	request := dtos.UpdateTodoDTO{
		Title:   "Hello World",
		Content: &updateTodoContent,
		Type:    enums.Text,
	}

	t.Logf("todo.ID before %v", todo.ID)
	todo.ID = todo.ID + 1
	t.Logf("todo.ID after %v", todo.ID)

	return UpdateTodoSetupResponse{
		User:       user,
		AuthToken:  authToken,
		Todo:       todo,
		RequestDto: request,
	}
}

func returnErrorOnIncorrectTodoIDExtraAssertions(t *testing.T, setup UpdateTodoSetupResponse) {
	t.Helper()
	todo := database.Todo{}
	result := TestServerInstance.DB.Where("id = ?", setup.Todo.ID).First(&todo)
	if result.Error == nil {
		t.Fatalf("todo should not exist in the db")
	}
}
