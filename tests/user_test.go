package tests

import (
	"encoding/json"
	"net/http"
	"path"
	"testing"

	utils_routes "github.com/Brawdunoir/dionysos-server/utils/routes"
	utils "github.com/Brawdunoir/dionysos-server/utils/tests"
)

// CreateResponseUser allows to map the response of the CreateUser request and get the key for further requests.
type CreateResponseUser utils_routes.CreateResponse

func (c CreateResponseUser) TargetURI(body []byte) (string, error) {
	err := json.Unmarshal(body, &c)

	return c.URI, err
}

func (c CreateResponseUser) Headers(body []byte) ([]utils.Header, error) {
	err := json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}

	return utils.GetBasicAuthHeader(path.Base(c.URI), c.Password), nil
}

var userURL = "/users"
var userCreateRequest = utils.Request{Method: http.MethodPost, Target: userURL, Body: `{"name":"test"}`}

// TestCreateUser tests the CreateUser function.
func TestCreateUser(t *testing.T) {
	method := http.MethodPost

	test := utils.TestCreate{
		Target: userURL,
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method, Body: `{"name":"test"}`}, ResponseCode: http.StatusCreated, ResponseBodyRegex: `{"uri":"` + userURL + `/\d+","password":"[0-9a-f]{64}"}`},
			{Name: "Empty body", Request: utils.Request{Method: method, Body: ``}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Empty json", Request: utils.Request{Method: method, Body: `{}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Bad name key", Request: utils.Request{Method: method, Body: `{"wrongkey":"test"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Empty name value", Request: utils.Request{Method: method, Body: `{"name":""}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Nil name value", Request: utils.Request{Method: method, Body: `{"name":nil}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Integer name value", Request: utils.Request{Method: method, Body: `{"name":1}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Object name value", Request: utils.Request{Method: method, Body: `{"name":{"somekey":"somevalue"}}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Less than min caracters", Request: utils.Request{Method: method, Body: `{"name":"a"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Exactly min caracters", Request: utils.Request{Method: method, Body: `{"name":"ab"}`}, ResponseCode: http.StatusCreated, ResponseBodyRegex: `{"uri":"` + userURL + `/\d+","password":"[0-9a-f]{64}"}`},
			{Name: "More than max caracters", Request: utils.Request{Method: method, Body: `{"name":"xxxxxxxxxxxxxxxxxxxxx"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Exactly max caracters", Request: utils.Request{Method: method, Body: `{"name":"xxxxxxxxxxxxxxxxxxxx"}`}, ResponseCode: http.StatusCreated, ResponseBodyRegex: `{"uri":"` + userURL + `/\d+","password":"[0-9a-f]{64}"}`},
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
			{Name: "Success", Request: utils.Request{Method: method}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"name":"test"}`},
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
			{Name: "Success", Request: utils.Request{Method: method, Body: `{"name":"test2"}`}, ResponseCode: http.StatusNoContent, ResponseBodyRegex: ``},
			{Name: "Correctly updated", Request: utils.Request{Method: http.MethodGet}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"name":"test2"}`},
			{Name: "Empty Body", Request: utils.Request{Method: method, Body: ``}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Empty json", Request: utils.Request{Method: method, Body: `{}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Bad name key", Request: utils.Request{Method: method, Body: `{"wrongkey":"test2"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Empty name value", Request: utils.Request{Method: method, Body: `{"name":""}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Nil name value", Request: utils.Request{Method: method, Body: `{"name":nil}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Integer name value", Request: utils.Request{Method: method, Body: `{"name":1}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Object name value", Request: utils.Request{Method: method, Body: `{"name":{"somekey":"somevalue"}}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Less than min caracters", Request: utils.Request{Method: method, Body: `{"name":"a"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Exactly min caracters", Request: utils.Request{Method: method, Body: `{"name":"ab"}`}, ResponseCode: http.StatusNoContent, ResponseBodyRegex: ``},
			{Name: "Correctly updated", Request: utils.Request{Method: http.MethodGet}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"name":"ab"}`},
			{Name: "More than max caracters", Request: utils.Request{Method: method, Body: `{"name":"xxxxxxxxxxxxxxxxxxxxx"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Exactly max caracters", Request: utils.Request{Method: method, Body: `{"name":"xxxxxxxxxxxxxxxxxxxx"}`}, ResponseCode: http.StatusNoContent, ResponseBodyRegex: ``},
			{Name: "Correctly updated", Request: utils.Request{Method: http.MethodGet}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"name":"xxxxxxxxxxxxxxxxxxxx"}`},
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
			{Name: "Invalid ID", Request: utils.Request{Method: method, Target: "abc"}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":"Invalid user ID"}`},
			{Name: "Success", Request: utils.Request{Method: method}, ResponseCode: http.StatusNoContent, ResponseBodyRegex: ``},
			{Name: "Correctly deleted", Request: utils.Request{Method: http.MethodGet}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"User not found"}`},
			{Name: "Not found", Request: utils.Request{Method: method, Target: "987654321"}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"User not found"}`},
		},
	}
	test.Run(t)
}
