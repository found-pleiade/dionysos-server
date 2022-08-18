package tests

import (
	"encoding/json"
	"fmt"
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
	err := database.MigrateDB(database.DB, true)
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
	err := database.MigrateDB(database.DB, true)
	if err != nil {
		t.Error(err)
	}

	// Create the user that will be used to pursue the tests.
	id, headers, err := utils.CreateTestUser(models.User{Name: "test"})
	if err != nil {
		t.Error(err)
	}

	suffix := fmt.Sprintf(`,"ownerID":%s,"users":\[{"ID":%s,"name":"test"}\]}`, id, id)

	method := http.MethodGet
	test := utils.TestRUD{
		CreateRequest:        roomCreateRequest,
		CreateResponse:       CreateResponseRoom{},
		CreateRequestHeaders: headers,
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method, Headers: headers}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"name":"test"` + suffix},
			{Name: "Invalid ID", Request: utils.Request{Method: method, Target: "abc", Headers: headers}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":"Invalid room ID"}`},
			{Name: "Not found", Request: utils.Request{Method: method, Target: "987654321", Headers: headers}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
		},
	}
	test.Run(t)
}

// TestUpdateRoom tests the UpdateRoom function.
func TestUpdateRoom(t *testing.T) {
	err := database.MigrateDB(database.DB, true)
	if err != nil {
		t.Error(err)
	}

	// Create the user that will be used to pursue the tests.
	id, headers, err := utils.CreateTestUser(models.User{Name: "test"})
	if err != nil {
		t.Error(err)
	}

	suffix := fmt.Sprintf(`,"ownerID":%s,"users":\[{"ID":%s,"name":"test"}\]}`, id, id)

	method := http.MethodPatch
	test := utils.TestRUD{
		CreateRequest:        roomCreateRequest,
		CreateRequestHeaders: headers,
		CreateResponse:       CreateResponseRoom{},
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method, Headers: headers, Body: `{"name":"test2"}`}, ResponseCode: http.StatusNoContent, ResponseBodyRegex: ``},
			{Name: "Correctly updated", Request: utils.Request{Method: http.MethodGet, Headers: headers}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"name":"test2"` + suffix},
			{Name: "Empty Body", Request: utils.Request{Method: method, Headers: headers, Body: ``}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Empty json", Request: utils.Request{Method: method, Headers: headers, Body: `{}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Bad name key", Request: utils.Request{Method: method, Headers: headers, Body: `{"wrongkey":"test2"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Empty name value", Request: utils.Request{Method: method, Headers: headers, Body: `{"name":""}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Nil name value", Request: utils.Request{Method: method, Headers: headers, Body: `{"name":nil}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Integer name value", Request: utils.Request{Method: method, Headers: headers, Body: `{"name":1}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Object name value", Request: utils.Request{Method: method, Headers: headers, Body: `{"name":{"somekey":"somevalue"}}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Less than min caracters", Request: utils.Request{Method: method, Headers: headers, Body: `{"name":"a"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Exactly min caracters", Request: utils.Request{Method: method, Headers: headers, Body: `{"name":"ab"}`}, ResponseCode: http.StatusNoContent, ResponseBodyRegex: ``},
			{Name: "Correctly updated", Request: utils.Request{Method: http.MethodGet, Headers: headers}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"name":"ab"` + suffix},
			{Name: "More than max caracters", Request: utils.Request{Method: method, Headers: headers, Body: `{"name":"xxxxxxxxxxxxxxxxxxxxx"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Exactly max caracters", Request: utils.Request{Method: method, Headers: headers, Body: `{"name":"xxxxxxxxxxxxxxxxxxxx"}`}, ResponseCode: http.StatusNoContent, ResponseBodyRegex: ``},
			{Name: "Correctly updated", Request: utils.Request{Method: http.MethodGet, Headers: headers}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"name":"xxxxxxxxxxxxxxxxxxxx"` + suffix},
			{Name: "Invalid ID", Request: utils.Request{Method: method, Headers: headers, Target: "abc"}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":"Invalid room ID"}`},
			{Name: "Not found", Request: utils.Request{Method: method, Headers: headers, Target: "987654321", Body: `{"name":"test2"}`}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
		},
	}
	test.Run(t)
}

// TestConnectRoom tests the ConnectUserToRoom function.
func TestConnectRoom(t *testing.T) {
	err := database.MigrateDB(database.DB, true)
	if err != nil {
		t.Error(err)
	}

	// Create the user that will be used to pursue the tests.
	id, headers, err := utils.CreateTestUser(models.User{Name: "test"})
	if err != nil {
		t.Error(err)
	}

	regex := fmt.Sprintf(`,"ownerID":%s,"users":\[{"ID":%s,"name":"test"}\]}`, id, id)

	method := http.MethodPatch
	target := "/connect"
	test := utils.TestRUD{
		CreateRequest:        roomCreateRequest,
		CreateRequestHeaders: headers,
		CreateResponse:       CreateResponseRoom{},
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: http.MethodGet, Headers: headers}, ResponseCode: http.StatusOK, ResponseBodyRegex: regex},
			{Name: "Connect 2nd time", Request: utils.Request{Target: target, Method: method, Headers: headers}, ResponseCode: http.StatusConflict, ResponseBodyRegex: `{"error":"User already in room"}`},
			{Name: "Not added 2nd time", Request: utils.Request{Method: http.MethodGet, Headers: headers}, ResponseCode: http.StatusOK, ResponseBodyRegex: regex},
			{Name: "Invalid ID", Request: utils.Request{Method: method, Headers: headers, Target: "abc" + target}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":"Invalid room ID"}`},
			{Name: "Not found", Request: utils.Request{Method: method, Headers: headers, Target: "987654321" + target}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
		},
	}
	test.Run(t)
}

// TestDisconnectRoom tests the DisconnectUserFromRoom function.
func TestDisconnectRoom(t *testing.T) {
	err := database.MigrateDB(database.DB, true)
	if err != nil {
		t.Error(err)
	}

	// Create the user that will be used to pursue the tests.
	_, headers, err := utils.CreateTestUser(models.User{Name: "test"})
	if err != nil {
		t.Error(err)
	}

	method := http.MethodPatch
	target := "/disconnect"
	test := utils.TestRUD{
		CreateRequest:        roomCreateRequest,
		CreateRequestHeaders: headers,
		CreateResponse:       CreateResponseRoom{},
		SubTests: []utils.SubTest{
			{Name: "Invalid ID", Request: utils.Request{Method: method, Headers: headers, Target: "abc" + target}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":"Invalid room ID"}`},
			{Name: "Not found", Request: utils.Request{Method: method, Headers: headers, Target: "987654321" + target}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
			{Name: "Success", Request: utils.Request{Target: target, Method: method, Headers: headers}, ResponseCode: http.StatusNoContent, ResponseBodyRegex: ``},
			{Name: "Room should be deleted", Request: utils.Request{Method: http.MethodGet, Headers: headers}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
		},
	}
	test.Run(t)
}

// TestRoomScenarioA is the following scenario:
// — 3 users (A, B, C) join the room, A is the owner.
// — C disconnects, nothing happens.
// — C tries to disconnect again, it is refused.
// — A disconnects, the ownership is transfered to B.
// — B disconnects, the room is deleted.
func TestRoomScenarioA(t *testing.T) {
	err := database.MigrateDB(database.DB, true)
	if err != nil {
		t.Error(err)
	}

	// Create the user that will be used to pursue the tests.
	idA, headersA, err := utils.CreateTestUser(models.User{Name: "userA"})
	if err != nil {
		t.Error(err)
	}
	idB, headersB, err := utils.CreateTestUser(models.User{Name: "userB"})
	if err != nil {
		t.Error(err)
	}
	idC, headersC, err := utils.CreateTestUser(models.User{Name: "userC"})
	if err != nil {
		t.Error(err)
	}

	name := `{"name":"test"`
	roomWhenA := fmt.Sprintf(`%s,"ownerID":%s,"users":\[{"ID":%s,"name":"userA"}\]}`, name, idA, idA)
	roomWhenAB := fmt.Sprintf(`%s,"ownerID":%s,"users":\[{"ID":%s,"name":"userA"},{"ID":%s,"name":"userB"}\]}`, name, idA, idA, idB)
	roomWhenABC := fmt.Sprintf(`%s,"ownerID":%s,"users":\[{"ID":%s,"name":"userA"},{"ID":%s,"name":"userB"},{"ID":%s,"name":"userC"}\]}`, name, idA, idA, idB, idC)
	roomWhenB := fmt.Sprintf(`%s,"ownerID":%s,"users":\[{"ID":%s,"name":"userB"}\]}`, name, idB, idB)

	targetConnect := "/connect"
	targetDisconnect := "/disconnect"

	test := utils.TestRUD{
		CreateRequest:        roomCreateRequest,
		CreateRequestHeaders: headersA,
		CreateResponse:       CreateResponseRoom{},
		SubTests: []utils.SubTest{
			{Name: "Assert A is in room", Request: utils.Request{Method: http.MethodGet, Headers: headersB}, ResponseCode: http.StatusOK, ResponseBodyRegex: roomWhenA},
			{Name: "B joins", Request: utils.Request{Target: targetConnect, Method: http.MethodPatch, Headers: headersB}, ResponseCode: http.StatusNoContent, ResponseBodyRegex: ``},
			{Name: "Assert B has joined", Request: utils.Request{Method: http.MethodGet, Headers: headersB}, ResponseCode: http.StatusOK, ResponseBodyRegex: roomWhenAB},
			{Name: "C joins", Request: utils.Request{Target: targetConnect, Method: http.MethodPatch, Headers: headersC}, ResponseCode: http.StatusNoContent, ResponseBodyRegex: ``},
			{Name: "Assert C has joined", Request: utils.Request{Method: http.MethodGet, Headers: headersC}, ResponseCode: http.StatusOK, ResponseBodyRegex: roomWhenABC},
			{Name: "C disconnects", Request: utils.Request{Target: targetDisconnect, Method: http.MethodPatch, Headers: headersC}, ResponseCode: http.StatusNoContent, ResponseBodyRegex: ``},
			{Name: "Assert C has disconnected", Request: utils.Request{Method: http.MethodGet, Headers: headersA}, ResponseCode: http.StatusOK, ResponseBodyRegex: roomWhenAB},
			{Name: "C tries to disconnects again", Request: utils.Request{Target: targetDisconnect, Method: http.MethodPatch, Headers: headersC}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":"User not in room"}`},
			{Name: "A disconnects", Request: utils.Request{Target: targetDisconnect, Method: http.MethodPatch, Headers: headersA}, ResponseCode: http.StatusNoContent, ResponseBodyRegex: ``},
			{Name: "Assert A has disconnected", Request: utils.Request{Method: http.MethodGet, Headers: headersA}, ResponseCode: http.StatusOK, ResponseBodyRegex: roomWhenB},
			{Name: "B disconnects", Request: utils.Request{Target: targetDisconnect, Method: http.MethodPatch, Headers: headersB}, ResponseCode: http.StatusNoContent, ResponseBodyRegex: ``},
			{Name: "Room should be deleted", Request: utils.Request{Method: http.MethodGet, Headers: headersA}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
		},
	}
	test.Run(t)
}
