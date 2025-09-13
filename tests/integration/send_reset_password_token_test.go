package integration

import (
	"bytes"
	"encoding/json"
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/horlerdipo/todo-golang/utils"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

var initiateResetPasswordRequest = struct {
	Email string `json:"email"`
}{
	Email: "john.doe@example.com",
}

func TestResetPasswordToken_ValidationError(t *testing.T) {
	//ARRANGE:
	ClearAllTables(t, TestServerInstance.DB)
	seedUser(t, initiateResetPasswordRequest)
	request, err := json.Marshal(struct{}{})
	if err != nil {
		t.Fatal("Unable to marshal initiate reset password request", err)
	}

	//ACT:
	response, err := http.Post(TestServerInstance.Server.URL+"/auth/password/forgot", "application/json", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal("Unable to send initiate reset password request", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %s", err)
	}

	var responseJson utils.JsonResponse[struct{}]
	err = json.Unmarshal(responseBody, &responseJson)

	//ASSERT:
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, responseJson.Message, "Validation error")
}

func TestResetPasswordToken_UnregisteredEmail(t *testing.T) {
	//ARRANGE
	ClearAllTables(t, TestServerInstance.DB)
	request, err := json.Marshal(initiateResetPasswordRequest)
	if err != nil {
		t.Fatal("Unable to marshal initiate reset password request", err)
	}

	//ACT
	response, err := http.Post(TestServerInstance.Server.URL+"/auth/password/forgot", "application/json", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal("Unable to send initiate reset password request", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %s", err)
	}

	var responseJson utils.JsonResponse[struct{}]
	err = json.Unmarshal(responseBody, &responseJson)

	//ASSERT
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, responseJson.Message, "email does not exist")
}

func TestResetPasswordToken_Success(t *testing.T) {
	//ARRANGE
	ClearAllTables(t, TestServerInstance.DB)
	seedUser(t, initiateResetPasswordRequest)
	request, err := json.Marshal(initiateResetPasswordRequest)
	if err != nil {
		t.Fatal("Unable to marshal initiate reset password request", err)
	}

	//ACT
	response, err := http.Post(TestServerInstance.Server.URL+"/auth/password/forgot", "application/json", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal("Unable to send initiate reset password request", err)
	}
	defer response.Body.Close()

	//ASSERT
	assert.Equal(t, http.StatusNoContent, response.StatusCode)
	user := database.User{}
	result := TestServerInstance.DB.First(&user, "email = ?", initiateResetPasswordRequest.Email)
	if result.Error != nil {
		t.Fatal("Unable to find user with email", initiateResetPasswordRequest.Email)
	}
	assert.NotEmpty(t, user.ResetToken)
	assert.NotEmpty(t, user.ResetTokenExpiresAt)
}
