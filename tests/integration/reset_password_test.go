package integration

import (
	"bytes"
	"encoding/json"
	"github.com/horlerdipo/todo-golang/internal/users"
	"github.com/horlerdipo/todo-golang/utils"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
	"time"
)

var resetPasswordRequest = struct {
	NewPassword string `json:"new_password"`
	Token       string `json:"reset_token"`
}{
	NewPassword: "new-password",
	Token:       "908447",
}

func TestResetPassword_ValidationError(t *testing.T) {
	//ARRANGE:
	ClearAllTables(t, TestServerInstance.DB)
	request, err := json.Marshal(struct{}{})
	if err != nil {
		t.Fatal("Unable to marshal reset password request", err)
	}

	//ACT:
	response, err := http.Post(TestServerInstance.Server.URL+"/auth/password/reset", "application/json", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal("Unable to send reset password request", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal("Failed to read response body", err)
	}

	var responseJson utils.JsonResponse[struct{}]
	err = json.Unmarshal(responseBody, &responseJson)

	//ASSERT:
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, responseJson.Message, "Validation error")
}

func TestResetPassword_IncorrectToken(t *testing.T) {
	//ARRANGE
	ClearAllTables(t, TestServerInstance.DB)
	resetTokenTime := time.Now().Add(time.Hour * 24)
	seedUser(t, struct {
		ResetToken          *string
		ResetTokenExpiresAt *time.Time
	}{
		ResetToken:          &resetPasswordRequest.Token,
		ResetTokenExpiresAt: &resetTokenTime,
	})
	resetPasswordRequest.Token = "898444774"
	request, err := json.Marshal(resetPasswordRequest)
	if err != nil {
		t.Fatal("Unable to marshal reset password request", err)
	}

	//ACT
	response, err := http.Post(TestServerInstance.Server.URL+"/auth/password/reset", "application/json", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal("Unable to send reset password request", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %s", err)
	}

	var responseJson utils.JsonResponse[struct{}]
	err = json.Unmarshal(responseBody, &responseJson)
	if err != nil {
		t.Fatal("Failed to read response body", err)
	}

	//ASSERT
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "reset token is invalid", responseJson.Message)
}

func TestResetPassword_ExpiredToken(t *testing.T) {
	//ARRANGE
	ClearAllTables(t, TestServerInstance.DB)
	resetTokenTime := time.Now().Add(-time.Hour * 24)
	seedUser(t, struct {
		ResetToken          *string
		ResetTokenExpiresAt *time.Time
	}{
		ResetToken:          &resetPasswordRequest.Token,
		ResetTokenExpiresAt: &resetTokenTime,
	})

	request, err := json.Marshal(resetPasswordRequest)
	if err != nil {
		t.Fatal("Unable to marshal reset password request", err)
	}

	//ACT
	response, err := http.Post(TestServerInstance.Server.URL+"/auth/password/reset", "application/json", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal("Unable to send reset password request", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %s", err)
	}

	var responseJson utils.JsonResponse[struct{}]
	err = json.Unmarshal(responseBody, &responseJson)
	if err != nil {
		t.Fatal("Failed to read response body", err)
	}

	//ASSERT
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "reset token has expired", responseJson.Message)
}

func TestResetPassword_Success(t *testing.T) {
	//ARRANGE
	ClearAllTables(t, TestServerInstance.DB)
	resetTokenTime := time.Now().Add(time.Hour * 24)
	user := seedUser(t, struct {
		ResetToken          *string
		ResetTokenExpiresAt *time.Time
	}{
		ResetToken:          &resetPasswordRequest.Token,
		ResetTokenExpiresAt: &resetTokenTime,
	})

	request, err := json.Marshal(resetPasswordRequest)
	if err != nil {
		t.Fatal("Unable to marshal reset password request", err)
	}

	//ACT
	response, err := http.Post(TestServerInstance.Server.URL+"/auth/password/reset", "application/json", bytes.NewBuffer(request))
	if err != nil {
		t.Fatal("Unable to send reset password request", err)
	}
	defer response.Body.Close()

	//ASSERT
	assert.Equal(t, http.StatusNoContent, response.StatusCode)
	newUser := users.User{}
	result := TestServerInstance.DB.First(&newUser, "id = ?", user.ID)
	if result.Error != nil {
		t.Fatal("User not found", result.Error)
	}

	assert.Nil(t, newUser.ResetToken)
	assert.Nil(t, newUser.ResetTokenExpiresAt)
	assert.True(t, utils.CheckPasswordHash(resetPasswordRequest.NewPassword, newUser.Password))
}
