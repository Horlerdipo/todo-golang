package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/horlerdipo/todo-golang/internal/enums"
	"github.com/horlerdipo/todo-golang/utils"
	"io"
	"net/http"
	"testing"
)

func SetupUpdateChecklistTest(t *testing.T) (*database.User, string, *database.Checklist) {
	t.Helper()

	ClearAllTables(t, TestServerInstance.DB)
	user := SeedUser(t, struct{}{})
	todo := SeedTodo(t, struct {
		Type enums.TodoType
	}{
		Type: enums.Checklist,
	}, user.ID)

	checklist := SeedChecklist(t, struct{}{}, todo.ID)
	authToken := GenerateTestJwtToken(t, user.ID)
	return user, authToken, checklist
}

type UpdateItemOnTodoChecklistSetupResponse struct {
	AuthToken     string
	User          *database.User
	Checklist     *database.Checklist
	ChecklistItem string
}

func TestUpdateChecklistOnTodo(t *testing.T) {
	tests := []struct {
		description        string
		setupFunc          func(t *testing.T) UpdateItemOnTodoChecklistSetupResponse
		expectedStatusCode int
		expectedMsg        string
		extraAssertions    func(t *testing.T, setupFuncResponse *UpdateItemOnTodoChecklistSetupResponse)
	}{
		{
			description:        "checklist can be updated successfully",
			setupFunc:          updateChecklistSuccessfullySetup,
			expectedStatusCode: http.StatusNoContent,
			expectedMsg:        "",
			extraAssertions:    updateChecklistSuccessfullyExtraAssertions,
		},
		{
			description:        "checklist returns validation error",
			setupFunc:          updateItemOnTodoChecklistValidationErrorSetup,
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedMsg:        "",
			extraAssertions:    updateItemOnTodoChecklistValidationErrorExtraAssertions,
		},
		{
			description:        "returns error when incorrect Todo ID is used",
			setupFunc:          updateItemOnToChecklistWithIncorrectTodoSetup,
			expectedStatusCode: http.StatusBadRequest,
			expectedMsg:        "todo does not exist",
			extraAssertions:    nil,
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

			url := fmt.Sprintf("%s/todos/%d/checklist/%d", TestServerInstance.Server.URL, setup.Checklist.TodoID, setup.Checklist.ID)
			req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonRequest))
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
				t.Log(string(response), resp.Request.URL.String())
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

				if tt.extraAssertions != nil {
					tt.extraAssertions(subTest, &setup)
				}
			}
		})
	}
}

func updateChecklistSuccessfullySetup(t *testing.T) UpdateItemOnTodoChecklistSetupResponse {
	t.Helper()
	user, authToken, checklist := SetupUpdateChecklistTest(t)
	return UpdateItemOnTodoChecklistSetupResponse{
		AuthToken:     authToken,
		User:          user,
		Checklist:     checklist,
		ChecklistItem: "Testing",
	}
}

func updateChecklistSuccessfullyExtraAssertions(t *testing.T, setup *UpdateItemOnTodoChecklistSetupResponse) {
	t.Helper()

	checklist := database.Checklist{}
	result := TestServerInstance.DB.Where("").First(&checklist)
	if result.Error != nil {
		t.Errorf("failed to find checklist: %v", result.Error)
	}

	if checklist.Description != setup.ChecklistItem {
		t.Errorf("expected checklist description %v, but got %v", setup.ChecklistItem, checklist.Description)
	}
}

func updateItemOnTodoChecklistValidationErrorSetup(t *testing.T) UpdateItemOnTodoChecklistSetupResponse {
	t.Helper()
	user, authToken, checklist := SetupUpdateChecklistTest(t)
	return UpdateItemOnTodoChecklistSetupResponse{
		AuthToken:     authToken,
		User:          user,
		Checklist:     checklist,
		ChecklistItem: "",
	}
}

func updateItemOnTodoChecklistValidationErrorExtraAssertions(t *testing.T, setup *UpdateItemOnTodoChecklistSetupResponse) {
	t.Helper()
	checklist := database.Checklist{}
	result := TestServerInstance.DB.Where("id", checklist.ID).First(&checklist)
	if result.Error != nil {
		t.Errorf("failed to find checklist: %v", result.Error)
	}

	if checklist.Description == setup.ChecklistItem {
		t.Errorf("expected checklist description %v, but got %v", checklist.Description, setup.ChecklistItem)
	}
}

func updateItemOnToChecklistWithIncorrectTodoSetup(t *testing.T) UpdateItemOnTodoChecklistSetupResponse {
	t.Helper()
	user, authToken, checklist := SetupUpdateChecklistTest(t)
	checklist.TodoID = checklist.TodoID + 1
	return UpdateItemOnTodoChecklistSetupResponse{
		AuthToken:     authToken,
		User:          user,
		Checklist:     checklist,
		ChecklistItem: "up manchester?",
	}
}
