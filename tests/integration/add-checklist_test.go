package integration

import (
	"bytes"
	"encoding/json"
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/horlerdipo/todo-golang/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func SetupAddChecklistTest(t *testing.T) (*database.User, string, *database.Todo) {
	ClearAllTables(t, TestServerInstance.DB)
	user := SeedUser(t, struct{}{})
	todo := SeedTodo(t, struct{}{}, user.ID)
	authToken := GenerateTestJwtToken(t, user.ID)
	return user, authToken, todo
}

func TestAddChecklist_Success(t *testing.T) {
	//ARRANGE
	_, authToken, todo := SetupPinTest(t)

	req, err := http.NewRequest(http.MethodPatch, TestServerInstance.Server.URL+"/todos/"+strconv.Itoa(int(todo.ID))+"/pin", bytes.NewBuffer(nil))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)
	//ACT
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	//ASSERT
	require.Equal(t, http.StatusOK, resp.StatusCode)
	newTodo := &database.Todo{}
	result := TestServerInstance.DB.Where("id = ?", todo.ID).First(&newTodo)
	assert.NoError(t, result.Error)
	assert.Equal(t, true, newTodo.Pinned)
}

func TestAddChecklist_NotFoundError(t *testing.T) {
	//ARRANGE
	ClearAllTables(t, TestServerInstance.DB)
	_, authToken, todo := SetupPinTest(t)

	req, err := http.NewRequest(http.MethodPatch, TestServerInstance.Server.URL+"/todos/"+strconv.Itoa(int(todo.ID-1))+"/pin", bytes.NewBuffer(nil))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	//ACT
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var responseJson utils.JsonResponse[interface{}]
	err = json.Unmarshal(body, &responseJson)
	require.NoError(t, err)

	//ASSERT
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, responseJson.Message, "todo does not exist")
}

func TestAddChecklist_UnauthorizedError(t *testing.T) {
	//ARRANGE
	ClearAllTables(t, TestServerInstance.DB)
	_, _, todo := SetupPinTest(t)

	req, err := http.NewRequest(http.MethodPatch, TestServerInstance.Server.URL+"/todos/"+strconv.Itoa(int(todo.ID-1))+"/pin", bytes.NewBuffer(nil))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	//ACT
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var responseJson utils.JsonResponse[interface{}]
	err = json.Unmarshal(body, &responseJson)
	require.NoError(t, err)

	//ASSERT
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Equal(t, responseJson.Message, "Unauthenticated")
}
