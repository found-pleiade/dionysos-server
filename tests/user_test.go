package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	utils "github.com/Brawdunoir/dionysos-server/utils/tests"
)

// CreateResponseUser allows to map the response of the CreateUser request and get the key for further requests.
type CreateResponseUser struct {
	ID string `json:"id"`
}

func (c CreateResponseUser) KeyCreated(body []byte) (string, error) {
	err := json.Unmarshal(body, &c)

	return c.ID, err
}

// TestCreateUser tests the CreateUser function.
func TestCreateUser(t *testing.T) {
	method, url := http.MethodPost, "/users/"

	test := utils.TestCreate{
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method, Url: url, Body: `{"username":"test"}`}, ResponseCode: http.StatusCreated, ResponseBodyRegex: `{"id":"\d+"}`},
			{Name: "Empty body", Request: utils.Request{Method: method, Url: url, Body: ``}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Bad username key", Request: utils.Request{Method: method, Url: url, Body: `{"wrongkey":"test"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Bad username value", Request: utils.Request{Method: method, Url: url, Body: `{"username":""}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
		},
	}
	test.Run(t)
}

// TestGetUser test the GetUser function.
func TestGetUser(t *testing.T) {
	method, url := http.MethodGet, "/users/"

	test := utils.TestGet{
		CreateRequest:  utils.Request{Method: http.MethodPost, Url: url, Body: `{"username":"test"}`},
		CreateResponse: CreateResponseUser{},
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method, Url: url}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"user":{.+}}`},
			{Name: "Not found", Request: utils.Request{Method: method, Url: url + "0"}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"User not found"}`},
		},
	}
	test.Run(t)
}
