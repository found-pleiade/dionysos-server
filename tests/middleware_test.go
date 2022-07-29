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
			{Name: "Wrong password", Request: utils.Request{Method: method}, ResponseCode: http.StatusUnauthorized, Headers: utils.GetBasicAuthHeader(id, "password"), ResponseBodyRegex: `{"error":"User not authorized"}`},
			{Name: "User not found", Request: utils.Request{Method: method}, ResponseCode: http.StatusNotFound, Headers: utils.GetBasicAuthHeader("987654321", "password"), ResponseBodyRegex: `{"error":"User not found"}`},
			{Name: "Empty authorization header", Request: utils.Request{Method: method}, ResponseCode: http.StatusUnauthorized, Headers: []utils.Header{}, ResponseBodyRegex: `{"error":"User not authorized"}`},
		},
	}

	tests.Run(t)
}
