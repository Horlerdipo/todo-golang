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

func SetupUpdateChecklistStatusTest(t *testing.T) (*database.User, string, *database.Checklist) {
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

type UpdateItemStatusOnTodoChecklistSetupResponse struct {
	AuthToken string
	User      *database.User
	Checklist *database.Checklist
	Status    bool
}

func TestUpdateChecklistStatusOnTodo(t *testing.T) {
	tests := []struct {
		description        string
		setupFunc          func(t *testing.T) UpdateItemStatusOnTodoChecklistSetupResponse
		expectedStatusCode int
		expectedMsg        string
		extraAssertions    func(t *testing.T, setupFuncResponse *UpdateItemStatusOnTodoChecklistSetupResponse)
	}{
		{
			description:        "checklist can be marked as done successfully",
			setupFunc:          markChecklistAsDoneSuccessfullySetup,
			expectedStatusCode: http.StatusNoContent,
			expectedMsg:        "",
			extraAssertions:    markChecklistAsDoneSuccessfullyExtraAssertions,
		},
		{
			description:        "returns error when incorrect Todo ID is used",
			setupFunc:          updateChecklistStatusWithIncorrectTodoSetup,
			expectedStatusCode: http.StatusBadRequest,
			expectedMsg:        "todo does not exist",
			extraAssertions:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(subTest *testing.T) {

			setup := tt.setupFunc(subTest)

			jsonRequest, err := json.Marshal(map[string]bool{
				"done": setup.Status,
			})
			if err != nil {
				t.Fatal("Failed to marshal request")
			}

			url := fmt.Sprintf("%s/todos/%d/checklist/%d", TestServerInstance.Server.URL, setup.Checklist.TodoID, setup.Checklist.ID)
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

func markChecklistAsDoneSuccessfullySetup(t *testing.T) UpdateItemStatusOnTodoChecklistSetupResponse {
	t.Helper()
	user, authToken, checklist := SetupUpdateChecklistTest(t)
	return UpdateItemStatusOnTodoChecklistSetupResponse{
		AuthToken: authToken,
		User:      user,
		Checklist: checklist,
		Status:    true,
	}
}

func markChecklistAsDoneSuccessfullyExtraAssertions(t *testing.T, setup *UpdateItemStatusOnTodoChecklistSetupResponse) {
	t.Helper()

	checklist := database.Checklist{}
	result := TestServerInstance.DB.Where("id = ?", setup.Checklist.ID).First(&checklist)
	if result.Error != nil {
		t.Errorf("failed to find checklist: %v", result.Error)
	}

	if checklist.Done != setup.Status {
		t.Errorf("expected checklist status to be %v, but got %v", setup.Status, checklist.Done)
	}
}

func markChecklistAsNotDoneSuccessfullySetup(t *testing.T) UpdateItemStatusOnTodoChecklistSetupResponse {
	t.Helper()

	user, authToken, checklist := SetupUpdateChecklistTest(t)

	result := TestServerInstance.DB.Model(&checklist).Update("done", true)
	if result.Error != nil {
		t.Fatalf("failed to mark checklist as Done: %v", result.Error)
	}

	return UpdateItemStatusOnTodoChecklistSetupResponse{
		AuthToken: authToken,
		User:      user,
		Checklist: checklist,
		Status:    false,
	}
}

func markChecklistAsNotDoneSuccessfullyExtraAssertions(t *testing.T, setup *UpdateItemStatusOnTodoChecklistSetupResponse) {
	t.Helper()

	checklist := database.Checklist{}
	result := TestServerInstance.DB.Where("id = ?", setup.Checklist.ID).First(&checklist)
	if result.Error != nil {
		t.Errorf("failed to find checklist: %v", result.Error)
	}

	if checklist.Done != setup.Status {
		t.Errorf("expected checklist status to be %v, but got %v", setup.Status, checklist.Done)
	}
}

func updateChecklistStatusWithIncorrectTodoSetup(t *testing.T) UpdateItemStatusOnTodoChecklistSetupResponse {
	t.Helper()
	user, authToken, checklist := SetupUpdateChecklistTest(t)
	checklist.TodoID = checklist.TodoID + 1
	return UpdateItemStatusOnTodoChecklistSetupResponse{
		AuthToken: authToken,
		User:      user,
		Checklist: checklist,
		Status:    true,
	}
}
