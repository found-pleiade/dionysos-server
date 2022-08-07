package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Brawdunoir/dionysos-server/database"
	"github.com/Brawdunoir/dionysos-server/models"
	utils_routes "github.com/Brawdunoir/dionysos-server/utils/routes"
	utils "github.com/Brawdunoir/dionysos-server/utils/tests"
)

// CreateResponseRoom allows to map the response of the CreateRoom request and get the key for further requests.
type CreateResponseRoom utils_routes.CreateResponse

func (c CreateResponseRoom) TargetURI(body []byte) (string, error) {
	err := json.Unmarshal(body, &c)

	return c.URI, err
}

func (c CreateResponseRoom) Headers(body []byte) ([]utils.Header, error) {
	return []utils.Header{}, nil
}

var roomURL = "/rooms"
var roomCreateRequest = utils.Request{Method: http.MethodPost, Target: roomURL, Body: `{"name":"test"}`}

// TestCreateRoom tests the CreateRoom function.
func TestCreateRoom(t *testing.T) {
	err := utils.ResetTable(database.DB, &models.User{}, &models.Room{})
	if err != nil {
		t.Error(err)
	}

	// Create the user that will be used to pursue the tests.
	_, headers, err := utils.CreateTestUser(models.User{Name: "test"})
	if err != nil {
		t.Error(err)
	}

	method := http.MethodPost
	test := utils.TestCreate{
		Target:  roomURL,
		Headers: headers,
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method, Body: `{"name":"test"}`}, ResponseCode: http.StatusCreated, ResponseBodyRegex: `{"uri":"` + roomURL + `/\d+"}`},
			{Name: "Empty body", Request: utils.Request{Method: method, Body: ``}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Empty json", Request: utils.Request{Method: method, Body: `{}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Bad name key", Request: utils.Request{Method: method, Body: `{"wrongkey":"test"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Empty name value", Request: utils.Request{Method: method, Body: `{"name":""}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Nil name value", Request: utils.Request{Method: method, Body: `{"name":nil}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Integer name value", Request: utils.Request{Method: method, Body: `{"name":1}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Object name value", Request: utils.Request{Method: method, Body: `{"name":{"somekey":"somevalue"}}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Less than min caracters", Request: utils.Request{Method: method, Body: `{"name":"a"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Exactly min caracters", Request: utils.Request{Method: method, Body: `{"name":"ab"}`}, ResponseCode: http.StatusCreated, ResponseBodyRegex: `{"uri":"` + roomURL + `/\d+"}`},
			{Name: "More than max caracters", Request: utils.Request{Method: method, Body: `{"name":"xxxxxxxxxxxxxxxxxxxxx"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Exactly max caracters", Request: utils.Request{Method: method, Body: `{"name":"xxxxxxxxxxxxxxxxxxxx"}`}, ResponseCode: http.StatusCreated, ResponseBodyRegex: `{"uri":"` + roomURL + `/\d+"}`},
		},
	}
	test.Run(t)
}

// TestGetRoom tests the GetRoom function.
func TestGetRoom(t *testing.T) {
	err := utils.ResetTable(database.DB, &models.User{}, &models.Room{})
	if err != nil {
		t.Error(err)
	}

	// Create the user that will be used to pursue the tests.
	_, headers, err := utils.CreateTestUser(models.User{Name: "test"})
	if err != nil {
		t.Error(err)
	}

	method := http.MethodGet
	test := utils.TestRUD{
		CreateRequest:        roomCreateRequest,
		CreateResponse:       CreateResponseRoom{},
		CreateRequestHeaders: headers,
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method, Headers: headers}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"name":"test"}`},
			{Name: "Invalid ID", Request: utils.Request{Method: method, Target: "abc", Headers: headers}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":"Invalid room ID"}`},
			{Name: "Not found", Request: utils.Request{Method: method, Target: "987654321", Headers: headers}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
		},
	}
	test.Run(t)
}

// TestUpdateRoom tests the UpdateRoom function.
func TestUpdateRoom(t *testing.T) {
	err := utils.ResetTable(database.DB, &models.User{}, &models.Room{})
	if err != nil {
		t.Error(err)
	}

	// Create the user that will be used to pursue the tests.
	_, headers, err := utils.CreateTestUser(models.User{Name: "test"})
	if err != nil {
		t.Error(err)
	}

	method := http.MethodPatch
	test := utils.TestRUD{
		CreateRequest:        roomCreateRequest,
		CreateRequestHeaders: headers,
		CreateResponse:       CreateResponseRoom{},
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method, Headers: headers, Body: `{"name":"test2"}`}, ResponseCode: http.StatusNoContent, ResponseBodyRegex: ``},
			{Name: "Correctly updated", Request: utils.Request{Method: http.MethodGet, Headers: headers}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"name":"test2"}`},
			{Name: "Empty Body", Request: utils.Request{Method: method, Headers: headers, Body: ``}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Empty json", Request: utils.Request{Method: method, Headers: headers, Body: `{}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Bad name key", Request: utils.Request{Method: method, Headers: headers, Body: `{"wrongkey":"test2"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Empty name value", Request: utils.Request{Method: method, Headers: headers, Body: `{"name":""}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Nil name value", Request: utils.Request{Method: method, Headers: headers, Body: `{"name":nil}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Integer name value", Request: utils.Request{Method: method, Headers: headers, Body: `{"name":1}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Object name value", Request: utils.Request{Method: method, Headers: headers, Body: `{"name":{"somekey":"somevalue"}}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Less than min caracters", Request: utils.Request{Method: method, Headers: headers, Body: `{"name":"a"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Exactly min caracters", Request: utils.Request{Method: method, Headers: headers, Body: `{"name":"ab"}`}, ResponseCode: http.StatusNoContent, ResponseBodyRegex: ``},
			{Name: "Correctly updated", Request: utils.Request{Method: http.MethodGet, Headers: headers}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"name":"ab"}`},
			{Name: "More than max caracters", Request: utils.Request{Method: method, Headers: headers, Body: `{"name":"xxxxxxxxxxxxxxxxxxxxx"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Exactly max caracters", Request: utils.Request{Method: method, Headers: headers, Body: `{"name":"xxxxxxxxxxxxxxxxxxxx"}`}, ResponseCode: http.StatusNoContent, ResponseBodyRegex: ``},
			{Name: "Correctly updated", Request: utils.Request{Method: http.MethodGet, Headers: headers}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"name":"xxxxxxxxxxxxxxxxxxxx"}`},
			{Name: "Invalid ID", Request: utils.Request{Method: method, Headers: headers, Target: "abc"}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":"Invalid room ID"}`},
			{Name: "Not found", Request: utils.Request{Method: method, Headers: headers, Target: "987654321", Body: `{"name":"test2"}`}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
		},
	}
	test.Run(t)
}

// TestDeleteRoom tests the DeleteRoom function.
func TestDeleteRoom(t *testing.T) {
	err := utils.ResetTable(database.DB, &models.User{}, &models.Room{})
	if err != nil {
		t.Error(err)
	}

	// Create the user that will be used to pursue the tests.
	_, headers, err := utils.CreateTestUser(models.User{Name: "test"})
	if err != nil {
		t.Error(err)
	}

	method := http.MethodDelete
	test := utils.TestRUD{
		CreateRequest:        roomCreateRequest,
		CreateRequestHeaders: headers,
		CreateResponse:       CreateResponseRoom{},
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method, Headers: headers}, ResponseCode: http.StatusNoContent, ResponseBodyRegex: ``},
			{Name: "Correctly deleted", Request: utils.Request{Method: http.MethodGet, Headers: headers}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
			{Name: "Invalid ID", Request: utils.Request{Method: method, Headers: headers, Target: "abc"}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":"Invalid room ID"}`},
			{Name: "Not found", Request: utils.Request{Method: method, Headers: headers, Target: "987654321"}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
		},
	}
	test.Run(t)
}

// TestConnectUserToRoom tests the ConnectUserToRoom function.
func TestConnectUserToRoom(t *testing.T) {
	// Create the user that will be used to pursue the tests.
	_, headers, err := utils.CreateTestUser(models.User{Name: "test"})
	if err != nil {
		t.Error(err)
	}

	method := http.MethodPost
	test := utils.TestRUD{
		CreateRequest:        roomCreateRequest,
		CreateRequestHeaders: headers,
		CreateResponse:       CreateResponseRoom{},
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method, Headers: headers}, ResponseCode: http.StatusOK, ResponseBodyRegex: ``},
			{Name: "Invalid ID", Request: utils.Request{Method: method, Headers: headers, Target: "abc"}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":"Invalid room ID"}`},
			{Name: "Not found", Request: utils.Request{Method: method, Headers: headers, Target: "987654321"}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
		},
	}
	test.Run(t)
}

// TestDisconnectUserFromRoom tests the DisonnectUserFromRoom function.
func TestDiconnectUserFromRoom(t *testing.T) {
	// Create the user that will be used to pursue the tests.
	_, headers, err := utils.CreateTestUser(models.User{Name: "test"})
	if err != nil {
		t.Error(err)
	}

	method := http.MethodPost
	test := utils.TestRUD{
		CreateRequest:        roomCreateRequest,
		CreateRequestHeaders: headers,
		CreateResponse:       CreateResponseRoom{},
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method, Headers: headers}, ResponseCode: http.StatusOK, ResponseBodyRegex: ``},
			{Name: "Invalid ID", Request: utils.Request{Method: method, Headers: headers, Target: "abc"}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":"Invalid room ID"}`},
			{Name: "Not found", Request: utils.Request{Method: method, Headers: headers, Target: "987654321"}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
		},
	}
	test.Run(t)
}
