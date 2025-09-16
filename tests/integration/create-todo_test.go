package integration

import (
	"bytes"
	"encoding/json"
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"github.com/horlerdipo/todo-golang/internal/enums"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

var content string = "random random string"

var createTodoRequest dtos.CreateTodoDTO = dtos.CreateTodoDTO{
	Title:   "Hello World",
	Content: &content,
	Type:    enums.Text,
}

func setupTest(t *testing.T) (*database.User, string) {
	ClearAllTables(t, TestServerInstance.DB)
	user := SeedUser(t, struct{}{})
	authToken := GenerateTestJwtToken(t, user.ID)
	return user, authToken
}

func TestCreateTodo_SuccessOnText(t *testing.T) {
	//ARRANGE:
	user, authToken := setupTest(t)
	jsonRequest, err := json.Marshal(createTodoRequest)
	if err != nil {
		require.NoError(t, err, "Failed to marshal request")
	}

	httpRequest, err := http.NewRequest(http.MethodPost, TestServerInstance.Server.URL+"/todos", bytes.NewBuffer(jsonRequest))
	if err != nil {
		require.NoError(t, err, "Failed to create request")
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer "+authToken)

	//ACT:
	client := &http.Client{}
	response, err := client.Do(httpRequest)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	//ASSERT:
	assert.Equal(t, http.StatusCreated, response.StatusCode)
	todo := &database.Todo{}
	result := TestServerInstance.DB.Where("user_id = ?", user.ID).First(todo)
	assert.Nil(t, result.Error)
	assert.Equal(t, createTodoRequest.Title, todo.Title)
	assert.Equal(t, createTodoRequest.Type, todo.Type)
}

func TestCreateTodo_SuccessOnChecklist(t *testing.T) {
	//ARRANGE:
	user, authToken := setupTest(t)

	createTodoRequest.Type = enums.Checklist
	createTodoRequest.Checklist = []string{"checklist1", "checklist2", "checklist3"}
	jsonRequest, err := json.Marshal(createTodoRequest)
	if err != nil {
		require.NoError(t, err, "Failed to marshal request")
	}

	httpRequest, err := http.NewRequest(http.MethodPost, TestServerInstance.Server.URL+"/todos", bytes.NewBuffer(jsonRequest))
	if err != nil {
		require.NoError(t, err, "Failed to create request")
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer "+authToken)

	//ACT:
	client := &http.Client{}
	response, err := client.Do(httpRequest)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	//ASSERT:
	assert.Equal(t, http.StatusCreated, response.StatusCode)

	todo := &database.Todo{}
	result := TestServerInstance.DB.Where("user_id = ?", user.ID).First(todo)
	assert.Nil(t, result.Error)
	assert.Equal(t, createTodoRequest.Title, todo.Title)
	assert.Equal(t, createTodoRequest.Type, todo.Type)

	checklist := &database.Checklist{}
	result = TestServerInstance.DB.Where("todo_id = ?", todo.ID).First(checklist)
	assert.Nil(t, result.Error)
	assert.Contains(t, createTodoRequest.Checklist, checklist.Description)
	assert.Equal(t, false, checklist.Done)
}

func TestCreateTodo_ValidationError(t *testing.T) {
	//ARRANGE:
	_, authToken := setupTest(t)
	createTodoRequest.Title = ""
	jsonRequest, err := json.Marshal(createTodoRequest)
	if err != nil {
		require.NoError(t, err, "Failed to marshal request")
	}

	httpRequest, err := http.NewRequest(http.MethodPost, TestServerInstance.Server.URL+"/todos", bytes.NewBuffer(jsonRequest))
	if err != nil {
		require.NoError(t, err, "Failed to create request")
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer "+authToken)

	//ACT:
	client := &http.Client{}
	response, err := client.Do(httpRequest)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	//ASSERT:
	assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
}

func TestCreateTodo_UnauthorizedError(t *testing.T) {
	//ARRANGE:
	createTodoRequest.Title = ""
	jsonRequest, err := json.Marshal(createTodoRequest)
	if err != nil {
		require.NoError(t, err, "Failed to marshal request")
	}

	httpRequest, err := http.NewRequest(http.MethodPost, TestServerInstance.Server.URL+"/todos", bytes.NewBuffer(jsonRequest))
	if err != nil {
		require.NoError(t, err, "Failed to create request")
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	//ACT:
	client := &http.Client{}
	response, err := client.Do(httpRequest)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	//ASSERT:
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestCreateTodo_ValidationErrorOnChecklist(t *testing.T) {
	//ARRANGE:
	_, authToken := setupTest(t)

	createTodoRequest.Type = enums.Checklist
	jsonRequest, err := json.Marshal(createTodoRequest)
	if err != nil {
		require.NoError(t, err, "Failed to marshal request")
	}

	httpRequest, err := http.NewRequest(http.MethodPost, TestServerInstance.Server.URL+"/todos", bytes.NewBuffer(jsonRequest))
	if err != nil {
		require.NoError(t, err, "Failed to create request")
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer "+authToken)
	//ACT:
	client := &http.Client{}
	response, err := client.Do(httpRequest)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	//ASSERT:
	assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
}
