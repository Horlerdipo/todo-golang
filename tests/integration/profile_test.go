package integration

import (
	"encoding/json"
	"github.com/horlerdipo/todo-golang/internal/dtos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

func TestProfile_Success(t *testing.T) {
	//ARRANGE:
	ClearAllTables(t, TestServerInstance.DB)
	user := SeedUser(t, struct{}{})

	req, err := http.NewRequest("GET", TestServerInstance.Server.URL+"/auth/user", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+GenerateTestJwtToken(t, user.ID))
	req.Header.Set("Content-Type", "application/json")

	//ACT:
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal("unable to read response body", err)
	}

	var responseData struct {
		Data struct {
			dtos.UserDetailsDto
		}
	}
	err = json.Unmarshal(responseBody, &responseData)
	require.NoError(t, err, "unable to unmarshal response body")

	//ASSERT:
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, user.Email, responseData.Data.Email)
	assert.Equal(t, user.FirstName, responseData.Data.FirstName)
	assert.Equal(t, user.LastName, responseData.Data.LastName)
	assert.Equal(t, user.ID, responseData.Data.ID)
}

func TestProfile_Unauthorized(t *testing.T) {
	//ARRANGE:
	ClearAllTables(t, TestServerInstance.DB)

	//ACT:
	response, err := http.Get(TestServerInstance.Server.URL + "/auth/user")
	if err != nil {
		t.Fatal("unable to send request to server", err)
	}
	defer response.Body.Close()

	//ASSERT:
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestProfile_WrongAuthenticationToken(t *testing.T) {
	//ARRANGE:
	ClearAllTables(t, TestServerInstance.DB)
	SeedUser(t, struct{}{})

	req, err := http.NewRequest("GET", TestServerInstance.Server.URL+"/auth/user", nil)
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
