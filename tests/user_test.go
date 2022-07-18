package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Brawdunoir/dionysos-server/constants"
	utils "github.com/Brawdunoir/dionysos-server/utils/tests"
)

// CreateResponseUser allows to map the response of the CreateUser request and get the key for further requests.
type CreateResponseUser struct {
	URI string `json:"uri"`
}

func (c CreateResponseUser) TargetURI(body []byte) (string, error) {
	err := json.Unmarshal(body, &c)

	return c.URI, err
}

var userURL = constants.BasePath + "/users"
var userCreateRequest = utils.Request{Method: http.MethodPost, Target: userURL, Body: `{"name":"test"}`}

// TestCreateUser tests the CreateUser function.
func TestCreateUser(t *testing.T) {
	method := http.MethodPost

	test := utils.TestCreate{
		Target: userURL,
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method, Body: `{"name":"test"}`}, ResponseCode: http.StatusCreated, ResponseBodyRegex: `{"uri":"` + userURL + `/\d+"}`},
			{Name: "Empty body", Request: utils.Request{Method: method, Body: ``}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Empty json", Request: utils.Request{Method: method, Body: `{}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Bad name key", Request: utils.Request{Method: method, Body: `{"wrongkey":"test"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Empty name value", Request: utils.Request{Method: method, Body: `{"name":""}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Nil name value", Request: utils.Request{Method: method, Body: `{"name":nil}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Integer name value", Request: utils.Request{Method: method, Body: `{"name":1}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Object name value", Request: utils.Request{Method: method, Body: `{"name":{"somekey":"somevalue"}}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Less than min caracters", Request: utils.Request{Method: method, Body: `{"name":"a"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Exactly min caracters", Request: utils.Request{Method: method, Body: `{"name":"ab"}`}, ResponseCode: http.StatusCreated, ResponseBodyRegex: `{"uri":"` + userURL + `/\d+"}`},
			{Name: "More than max caracters", Request: utils.Request{Method: method, Body: `{"name":"xxxxxxxxxxxxxxxxxxxxx"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Exactly max caracters", Request: utils.Request{Method: method, Body: `{"name":"xxxxxxxxxxxxxxxxxxxx"}`}, ResponseCode: http.StatusCreated, ResponseBodyRegex: `{"uri":"` + userURL + `/\d+"}`},
		},
	}
	test.Run(t)
}

// TestGetUser tests the GetUser function.
func TestGetUser(t *testing.T) {
	method := http.MethodGet

	test := utils.TestRUD{
		CreateRequest:  userCreateRequest,
		CreateResponse: CreateResponseUser{},
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"user":{"name":"test"}}`},
			{Name: "Invalid ID", Request: utils.Request{Method: method, Target: "abc"}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":"Invalid user ID"}`},
			{Name: "Not found", Request: utils.Request{Method: method, Target: "987654321"}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"User not found"}`},
		},
	}
	test.Run(t)
}

// TestUpdateUser tests the UpdateUser function.
func TestUpdateUser(t *testing.T) {
	method := http.MethodPatch

	test := utils.TestRUD{
		CreateRequest:  userCreateRequest,
		CreateResponse: CreateResponseUser{},
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method, Body: `{"name":"test2"}`}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"user":{"name":"test2"}}`},
			{Name: "Correctly updated", Request: utils.Request{Method: http.MethodGet}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"user":{"name":"test2"}}`},
			{Name: "Empty Body", Request: utils.Request{Method: method, Body: ``}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Empty json", Request: utils.Request{Method: method, Body: `{}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Bad name key", Request: utils.Request{Method: method, Body: `{"wrongkey":"test2"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Empty name value", Request: utils.Request{Method: method, Body: `{"name":""}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Nil name value", Request: utils.Request{Method: method, Body: `{"name":nil}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Integer name value", Request: utils.Request{Method: method, Body: `{"name":1}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Object name value", Request: utils.Request{Method: method, Body: `{"name":{"somekey":"somevalue"}}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Less than min caracters", Request: utils.Request{Method: method, Body: `{"name":"a"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Exactly min caracters", Request: utils.Request{Method: method, Body: `{"name":"ab"}`}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"user":{"name":"ab"}}`},
			{Name: "More than max caracters", Request: utils.Request{Method: method, Body: `{"name":"xxxxxxxxxxxxxxxxxxxxx"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Exactly max caracters", Request: utils.Request{Method: method, Body: `{"name":"xxxxxxxxxxxxxxxxxxxx"}`}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"user":{"name":"xxxxxxxxxxxxxxxxxxxx"}}`},
			{Name: "Invalid ID", Request: utils.Request{Method: method, Target: "abc"}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":"Invalid user ID"}`},
			{Name: "Not found", Request: utils.Request{Method: method, Target: "987654321", Body: `{"name":"test2"}`}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"User not found"}`},
		},
	}
	test.Run(t)
}

// TestDeleteUser tests the DeleteUser function.
func TestDeleteUser(t *testing.T) {
	method := http.MethodDelete

	test := utils.TestRUD{
		CreateRequest:  userCreateRequest,
		CreateResponse: CreateResponseUser{},
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method}, ResponseCode: http.StatusOK, ResponseBodyRegex: ``},
			{Name: "Correctly deleted", Request: utils.Request{Method: http.MethodGet}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"User not found"}`},
			{Name: "Invalid ID", Request: utils.Request{Method: method, Target: "abc"}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":"Invalid user ID"}`},
			{Name: "Not found", Request: utils.Request{Method: method, Target: "987654321"}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"User not found"}`},
		},
	}
	test.Run(t)
}
