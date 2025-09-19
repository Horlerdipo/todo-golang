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

func SetupDeleteChecklistTest(t *testing.T) (*database.User, string, *database.Checklist) {
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

type DeleteItemFromChecklistSetupResponse struct {
	AuthToken string
	User      *database.User
	Checklist *database.Checklist
}

func TestRemoveItemFromChecklistToTodo(t *testing.T) {
	tests := []struct {
		description        string
		setupFunc          func(t *testing.T) DeleteItemFromChecklistSetupResponse
		expectedStatusCode int
		expectedMsg        string
		extraAssertions    func(t *testing.T, setupFuncResponse DeleteItemFromChecklistSetupResponse)
	}{
		{
			description:        "checklist item can be removed successfully",
			setupFunc:          removeChecklistItemSuccessfullySetup,
			expectedStatusCode: http.StatusNoContent,
			expectedMsg:        "",
			extraAssertions:    removeChecklistItemSuccessfullyExtraAssertions,
		},
		{
			description:        "returns unknown todo error",
			setupFunc:          unknownTodoSetup,
			expectedStatusCode: http.StatusBadRequest,
			expectedMsg:        "todo does not exist",
			extraAssertions:    unknownTodoExtraAssertions,
		},
		{
			description:        "returns unknown checklist error",
			setupFunc:          unknownChecklistSetup,
			expectedStatusCode: http.StatusNoContent,
			expectedMsg:        "",
			extraAssertions:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(subTest *testing.T) {

			setup := tt.setupFunc(subTest)

			url := fmt.Sprintf("%s/todos/%d/checklist/%d", TestServerInstance.Server.URL, setup.Checklist.TodoID, setup.Checklist.ID)
			req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(nil))
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
				var jsonResponse utils.JsonResponse[interface{}]
				err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
				if err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if jsonResponse.Message != tt.expectedMsg {
					t.Errorf("expected message %s, but got %s", tt.expectedMsg, jsonResponse.Message)
				}

				if tt.extraAssertions != nil {
					tt.extraAssertions(subTest, setup)
				}
			}
		})
	}
}

func removeChecklistItemSuccessfullySetup(t *testing.T) DeleteItemFromChecklistSetupResponse {
	user, authToken, checklist := SetupDeleteChecklistTest(t)
	return DeleteItemFromChecklistSetupResponse{
		User:      user,
		AuthToken: authToken,
		Checklist: checklist,
	}
}

func removeChecklistItemSuccessfullyExtraAssertions(t *testing.T, setup DeleteItemFromChecklistSetupResponse) {
	checklist := &database.Checklist{}
	result := TestServerInstance.DB.Where("id = ?", setup.Checklist.ID).First(&checklist)
	if result.Error == nil {
		t.Errorf("checklist item can't be removed")
	}
}

func unknownChecklistSetup(t *testing.T) DeleteItemFromChecklistSetupResponse {
	user, authToken, checklist := SetupDeleteChecklistTest(t)
	checklist.ID = checklist.ID + 1
	return DeleteItemFromChecklistSetupResponse{
		User:      user,
		AuthToken: authToken,
		Checklist: checklist,
	}
}

func unknownTodoSetup(t *testing.T) DeleteItemFromChecklistSetupResponse {
	user, authToken, checklist := SetupDeleteChecklistTest(t)
	checklist.Todo.ID = checklist.Todo.ID + 1
	checklist.TodoID = checklist.TodoID + 1

	return DeleteItemFromChecklistSetupResponse{
		User:      user,
		AuthToken: authToken,
		Checklist: checklist,
	}
}

func unknownTodoExtraAssertions(t *testing.T, setup DeleteItemFromChecklistSetupResponse) {
	checklist := &database.Checklist{}
	result := TestServerInstance.DB.Where("id = ?", setup.Checklist.ID).First(&checklist)
	if result.Error != nil {
		t.Errorf("expected error to be nil, but got %v", result.Error)
	}
}
