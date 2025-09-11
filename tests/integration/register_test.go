package integration

import (
	"bytes"
	"encoding/json"
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"github.com/horlerdipo/todo-golang/internal/users"
	"github.com/horlerdipo/todo-golang/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

var registerRequest = dtos.CreateUserDTO{
	FirstName: "John",
	LastName:  "Doe",
	Email:     "john.doe@example.com",
	Password:  "password123",
}

func TestRegister_ValidationError(t *testing.T) {
	//ARRANGE
	ClearAllTables(t, TestServerInstance.DB)
	request := map[string]string{}
	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal request")

	//ACT
	var resp *http.Response
	resp, err = http.Post(TestServerInstance.Server.URL+"/auth/register", "application/json", bytes.NewBuffer(jsonData))
	require.NoError(t, err, "Failed to make register request")
	defer resp.Body.Close()

	//ASSERT
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestRegister_EmailAlreadyRegistered(t *testing.T) {
	//ARRANGE
	ClearAllTables(t, TestServerInstance.DB)
	hashedPassword, _ := utils.HashPassword("password123")
	user := users.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     registerRequest.Email,
		Password:  hashedPassword,
	}

	result := TestServerInstance.DB.Create(&user)
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	jsonData, err := json.Marshal(registerRequest)
	require.NoError(t, err, "Failed to marshal request")

	//ACT
	resp, err := http.Post(TestServerInstance.Server.URL+"/auth/register", "application/json", bytes.NewBuffer(jsonData))
	require.NoError(t, err, "Failed to make register request")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %s", err)
	}
	var jsonResp utils.JsonResponse[interface{}]
	err = json.Unmarshal(body, &jsonResp)
	require.NoError(t, err, "Failed to marshal response")
	//ASSERT
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, jsonResp.Message, "user already exists")
}

func TestRegister_Success(t *testing.T) {

	//ARRANGE
	ClearAllTables(t, TestServerInstance.DB)
	jsonData, err := json.Marshal(registerRequest)
	require.NoError(t, err, "Failed to marshal register request")

	//ACT
	var resp *http.Response
	resp, err = http.Post(
		TestServerInstance.Server.URL+"/auth/register",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	require.NoError(t, err, "Failed to make register request")
	defer resp.Body.Close()

	// ASSERT:
	assert.Equal(t, http.StatusNoContent, resp.StatusCode, "Should return 204 for successful registration")

	var dbUser users.User
	result := TestServerInstance.DB.Where("email = ?", registerRequest.Email).First(&dbUser)
	assert.NoError(t, result.Error, "User should exist in database")
	assert.Equal(t, registerRequest.Email, dbUser.Email, "User should match email")
	assert.Equal(t, registerRequest.FirstName, dbUser.FirstName, "User should match first name")
	assert.Equal(t, registerRequest.LastName, dbUser.LastName, "User should match last name")
}
