package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Brawdunoir/dionysos-server/constants"
	utils "github.com/Brawdunoir/dionysos-server/utils/tests"
)

// CreateResponseRoom allows to map the response of the CreateRoom request and get the key for further requests.
type CreateResponseRoom struct {
	URI string `json:"uri"`
}

func (c CreateResponseRoom) TargetURI(body []byte) (string, error) {
	err := json.Unmarshal(body, &c)

	return c.URI, err
}

var roomURL = constants.BasePath + "/rooms"
var roomCreateRequest = utils.Request{Method: http.MethodPost, Target: roomURL, Body: `{"name":"test"}`}

// TestCreateRoom tests the CreateRoom function.
func TestCreateRoom(t *testing.T) {
	method := http.MethodPost

	test := utils.TestCreate{
		Target: roomURL,
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
	method := http.MethodGet

	test := utils.TestRUD{
		CreateRequest:  roomCreateRequest,
		CreateResponse: CreateResponseRoom{},
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"room":{"name":"test"}}`},
			{Name: "Invalid ID", Request: utils.Request{Method: method, Target: "abc"}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":"Invalid room ID"}`},
			{Name: "Not found", Request: utils.Request{Method: method, Target: "987654321"}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
		},
	}
	test.Run(t)
}

// TestUpdateRoom tests the UpdateRoom function.
func TestUpdateRoom(t *testing.T) {
	method := http.MethodPatch

	test := utils.TestRUD{
		CreateRequest:  roomCreateRequest,
		CreateResponse: CreateResponseRoom{},
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method, Body: `{"name":"test2"}`}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"room":{"name":"test2"}}`},
			{Name: "Correctly updated", Request: utils.Request{Method: http.MethodGet}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"room":{"name":"test2"}}`},
			{Name: "Empty Body", Request: utils.Request{Method: method, Body: ``}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Empty json", Request: utils.Request{Method: method, Body: `{}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Bad name key", Request: utils.Request{Method: method, Body: `{"wrongkey":"test2"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Empty name value", Request: utils.Request{Method: method, Body: `{"name":""}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Nil name value", Request: utils.Request{Method: method, Body: `{"name":nil}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Integer name value", Request: utils.Request{Method: method, Body: `{"name":1}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Object name value", Request: utils.Request{Method: method, Body: `{"name":{"somekey":"somevalue"}}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Less than min caracters", Request: utils.Request{Method: method, Body: `{"name":"a"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Exactly min caracters", Request: utils.Request{Method: method, Body: `{"name":"ab"}`}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"room":{"name":"ab"}}`},
			{Name: "More than max caracters", Request: utils.Request{Method: method, Body: `{"name":"xxxxxxxxxxxxxxxxxxxxx"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Exactly max caracters", Request: utils.Request{Method: method, Body: `{"name":"xxxxxxxxxxxxxxxxxxxx"}`}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"room":{"name":"xxxxxxxxxxxxxxxxxxxx"}}`},
			{Name: "Invalid ID", Request: utils.Request{Method: method, Target: "abc"}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":"Invalid room ID"}`},
			{Name: "Not found", Request: utils.Request{Method: method, Target: "987654321", Body: `{"name":"test2"}`}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
		},
	}
	test.Run(t)
}

// TestDeleteRoom tests the DeleteRoom function.
func TestDeleteRoom(t *testing.T) {
	method := http.MethodDelete

	test := utils.TestRUD{
		CreateRequest:  roomCreateRequest,
		CreateResponse: CreateResponseRoom{},
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method}, ResponseCode: http.StatusOK, ResponseBodyRegex: ``},
			{Name: "Correctly deleted", Request: utils.Request{Method: http.MethodGet}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
			{Name: "Invalid ID", Request: utils.Request{Method: method, Target: "abc"}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":"Invalid room ID"}`},
			{Name: "Not found", Request: utils.Request{Method: method, Target: "987654321"}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
		},
	}
	test.Run(t)
}
