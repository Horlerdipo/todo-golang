package integration

import (
	"bytes"
	"encoding/json"
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"github.com/horlerdipo/todo-golang/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

var loginRequest = dtos.LoginUserDTO{
	Email:    "johndoe@example.com",
	Password: "password123",
}

func TestLogin_ValidationError(t *testing.T) {
	//ARRANGE
	ClearAllTables(t, TestServerInstance.DB)
	request := map[string]string{}
	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal request")

	//ACT
	var resp *http.Response
	resp, err = http.Post(TestServerInstance.Server.URL+"/auth/login", "application/json", bytes.NewBuffer(jsonData))
	require.NoError(t, err, "Failed to make register request")
	defer resp.Body.Close()

	//ASSERT
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected 400 Bad Request")
}

func TestLogin_EmailDoesNotExist(t *testing.T) {
	//ARRANGE
	ClearAllTables(t, TestServerInstance.DB)
	seedUser(t, struct{}{})

	jsonData, err := json.Marshal(dtos.LoginUserDTO{
		Email:    "johndoe@kkdkd.com",
		Password: "password123",
	})
	require.NoError(t, err, "Failed to marshal request")

	//ACT
	resp, err := http.Post(TestServerInstance.Server.URL+"/auth/login", "application/json", bytes.NewBuffer(jsonData))
	require.NoError(t, err, "Failed to make login request")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %s", err)
	}
	var jsonResp utils.JsonResponse[dtos.LoginUserResponseDto]
	err = json.Unmarshal(body, &jsonResp)
	require.NoError(t, err, "Failed to marshal response")

	//ASSERT
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, jsonResp.Message, "email or password is not valid")
}

func TestLogin_WrongEmail(t *testing.T) {

	//ARRANGE
	ClearAllTables(t, TestServerInstance.DB)
	seedUser(t, struct{}{})
	request := loginRequest
	request.Email = "wrongEmail@gmail.com"

	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal login request")

	//ACT
	var resp *http.Response
	resp, err = http.Post(
		TestServerInstance.Server.URL+"/auth/login",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	require.NoError(t, err, "Failed to make login request")
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %s", err)
	}

	var jsonResp utils.JsonResponse[dtos.LoginUserResponseDto]
	err = json.Unmarshal(responseBody, &jsonResp)
	require.NoError(t, err, "Failed to marshal response")

	// ASSERT:
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400 for invalid email or password")
	assert.Equal(t, jsonResp.Message, "email or password is not valid")
}

func TestLogin_PasswordIncorrect(t *testing.T) {
	//ARRANGE:
	ClearAllTables(t, TestServerInstance.DB)
	seedUser(t, struct{}{})
	request := loginRequest
	request.Password = "wrongPassword"
	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal request")

	//ACT:
	response, err := http.Post(TestServerInstance.Server.URL+"/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to make login request: %s", err)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %s", err)
	}

	var jsonResp utils.JsonResponse[interface{}]
	err = json.Unmarshal(responseBody, &jsonResp)
	require.NoError(t, err, "Failed to marshal response")

	//ASSERT:
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, jsonResp.Message, "email or password is not valid")
}

func TestLogin_Success(t *testing.T) {
	//ARRANGE
	ClearAllTables(t, TestServerInstance.DB)
	seedUser[dtos.LoginUserDTO](t, loginRequest)
	request, err := json.Marshal(loginRequest)
	if err != nil {
		t.Fatalf("Failed to marshal login request: %s", err)
	}

	//ACT
	response, err := http.Post(TestServerInstance.Server.URL+"/auth/login", "application/json", bytes.NewBuffer(request))
	if err != nil {
		t.Fatalf("Failed to make login request: %s", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %s", err)
	}

	var responseJson utils.JsonResponse[dtos.LoginUserResponseDto]
	err = json.Unmarshal(responseBody, &responseJson)
	require.NoError(t, err, "Failed to marshal response")
	//panic(fmt.Sprintf("unreachable: +%v", responseJson))
	//ASSERT
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, loginRequest.Email, responseJson.Data.Email)
	assert.NotEmpty(t, responseJson.Data.Token)
}
