package tests

import (
	"net/http"
	"testing"

	"github.com/Brawdunoir/dionysos-server/models"
	utils "github.com/Brawdunoir/dionysos-server/utils/tests"
)

// TestAuthentication tests the authentication middleware.
// We test it on a single endpoint because behavior is the same for all endpoints.
func TestAuthenticate(t *testing.T) {
	id, headers, err := utils.CreateTestUser(models.User{Name: "test"})
	if err != nil {
		t.Error(err)
	}

	method := http.MethodPost
	tests := utils.TestCreate{
		Target:  roomURL,
		Headers: headers,
		SubTests: []utils.SubTest{
			{Name: "Wrong password", Request: utils.Request{Method: method, Headers: utils.GetBasicAuthHeader(id, "password")}, ResponseCode: http.StatusUnauthorized, ResponseBodyRegex: `{"error":"User not authorized"}`},
			{Name: "User not found", Request: utils.Request{Method: method, Headers: utils.GetBasicAuthHeader("987654321", "password")}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"User not found"}`},
			{Name: "Empty authorization header", Request: utils.Request{Method: method, Headers: []utils.Header{}}, ResponseCode: http.StatusUnauthorized, ResponseBodyRegex: `{"error":"User not authorized"}`},
			{Name: "Not well formed authorization header", Request: utils.Request{Method: method, Headers: []utils.Header{{Key: "Authorization", Value: "apikey xxx"}}}, ResponseCode: http.StatusUnauthorized, ResponseBodyRegex: `{"error":"User not authorized"}`},
		},
	}

	tests.Run(t)
}
