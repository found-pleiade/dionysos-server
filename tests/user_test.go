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

var url = "/users/"
var createRequest = utils.Request{Method: http.MethodPost, Url: url, Body: `{"username":"test"}`}

// TestCreateUser tests the CreateUser function.
func TestCreateUser(t *testing.T) {
	method := http.MethodPost

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
	method := http.MethodGet

	test := utils.TestRUD{
		CreateRequest:  createRequest,
		CreateResponse: CreateResponseUser{},
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method, Url: url}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"user":{.+}}`},
			{Name: "Not found", Request: utils.Request{Method: method, Url: url + "0"}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"User not found"}`},
		},
	}
	test.Run(t)
}

// TestUpdateUser test the UpdateUser function.
func TestUpdateUser(t *testing.T) {
	method := http.MethodPatch

	test := utils.TestRUD{
		CreateRequest:  createRequest,
		CreateResponse: CreateResponseUser{},
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method, Url: url, Body: `{"username":"test2"}`}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"user":{.+}}`},
			{Name: "Empty Body", Request: utils.Request{Method: method, Url: url, Body: ``}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Bad username key", Request: utils.Request{Method: method, Url: url, Body: `{"wrongkey":"test2"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Bad username value", Request: utils.Request{Method: method, Url: url, Body: `{"username":""}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Not found", Request: utils.Request{Method: method, Url: url + "0", Body: `{"username":"test2"}`}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"User not found"}`},
		},
	}
	test.Run(t)
}
