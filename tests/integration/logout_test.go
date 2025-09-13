package integration

import (
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestLogout_Success(t *testing.T) {
	//ARRANGE:
	ClearAllTables(t, TestServerInstance.DB)
	user := seedUser(t, struct{}{})

	req, err := http.NewRequest("POST", TestServerInstance.Server.URL+"/auth/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	authToken := GenerateTestJwtToken(t, user.ID)
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	//ACT:
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	//ASSERT:
	assert.Equal(t, http.StatusNoContent, response.StatusCode)
	blacklistedToken := database.TokenBlacklist{}
	result := TestServerInstance.DB.Where("token = ?", authToken).First(&blacklistedToken)
	assert.NoError(t, result.Error)
}

func TestLogout_Unauthorized(t *testing.T) {
	//ARRANGE:
	ClearAllTables(t, TestServerInstance.DB)

	//ACT:
	response, err := http.Post(TestServerInstance.Server.URL+"/auth/logout", "application/json", nil)
	if err != nil {
		t.Fatal("unable to send request to server", err)
	}
	defer response.Body.Close()

	//ASSERT:
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestLogout_WrongAuthenticationToken(t *testing.T) {
	//ARRANGE:
	ClearAllTables(t, TestServerInstance.DB)
	seedUser(t, struct{}{})

	req, err := http.NewRequest("POST", TestServerInstance.Server.URL+"/auth/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer random-string-that-is-definitely-definitely-not-a-correct-authentication-token")
	req.Header.Set("Content-Type", "application/json")

	//ACT:
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	//ASSERT:
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestLogout_BlacklistedToken(t *testing.T) {
	//ARRANGE:
	ClearAllTables(t, TestServerInstance.DB)
	user := seedUser(t, struct{}{})

	req, err := http.NewRequest("POST", TestServerInstance.Server.URL+"/auth/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	authToken := GenerateTestJwtToken(t, user.ID)
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	//add the token to blacklist table before
	now := time.Now()
	tokenBlacklist := database.TokenBlacklist{
		Token:     authToken,
		ExpiresAt: &now,
	}

	result := TestServerInstance.DB.Create(&tokenBlacklist)
	if result.Error != nil {
		t.Fatal("unable to add token to blacklist: ", result.Error)
	}

	//ACT:
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	//ASSERT:
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}
