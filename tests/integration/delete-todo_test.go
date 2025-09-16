package integration

import (
	"bytes"
	"encoding/json"
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/horlerdipo/todo-golang/internal/enums"
	"github.com/horlerdipo/todo-golang/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func SetupTest(t *testing.T) (*database.User, string) {
	ClearAllTables(t, TestServerInstance.DB)
	user := SeedUser(t, struct{}{})
	authToken := GenerateTestJwtToken(t, user.ID)
	return user, authToken
}

func TestDeleteTodo_Success(t *testing.T) {
	//ARRANGE
	ClearAllTables(t, TestServerInstance.DB)
	user, authToken := SetupTest(t)

	content := "test"
	todo := &database.Todo{
		Title:   "Hello",
		Content: &content,
		Type:    enums.Text,
		UserID:  user.ID,
	}

	result := TestServerInstance.DB.Create(&todo)
	require.NoError(t, result.Error)

	req, err := http.NewRequest(http.MethodDelete, TestServerInstance.Server.URL+"/todos/"+strconv.Itoa(int(todo.ID)), bytes.NewBuffer(nil))
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
	result = TestServerInstance.DB.Where("id = ?", todo.ID).First(&todo)
	assert.Error(t, result.Error)
}

func TestDeleteTodo_NotFoundError(t *testing.T) {
	//ARRANGE
	ClearAllTables(t, TestServerInstance.DB)
	_, authToken := SetupTest(t)

	req, err := http.NewRequest(http.MethodDelete, TestServerInstance.Server.URL+"/todos/"+strconv.Itoa(10), bytes.NewBuffer(nil))
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
