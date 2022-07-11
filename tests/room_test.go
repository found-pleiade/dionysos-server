package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	utils "github.com/Brawdunoir/dionysos-server/utils/tests"
)

// CreateResponseRoom allows to map the response of the CreateRoom request and get the key for further requests.
type CreateResponseRoom struct {
	ID string `json:"id"`
}

func (c CreateResponseRoom) KeyCreated(body []byte) (string, error) {
	err := json.Unmarshal(body, &c)

	return c.ID, err
}

var roomURL = "/rooms/"
var roomCreateRequest = utils.Request{Method: http.MethodPost, URL: roomURL, Body: `{"name":"test"}`}

// TestCreateRoom tests the CreateRoom function.
func TestCreateRoom(t *testing.T) {
	method := http.MethodPost

	test := utils.TestCreate{
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method, URL: roomURL, Body: `{"name":"test"}`}, ResponseCode: http.StatusCreated, ResponseBodyRegex: `{"id":"\d+"}`},
			{Name: "Empty body", Request: utils.Request{Method: method, URL: roomURL, Body: ``}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Empty json", Request: utils.Request{Method: method, URL: roomURL, Body: `{}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Bad name key", Request: utils.Request{Method: method, URL: roomURL, Body: `{"wrongkey":"test"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Empty name value", Request: utils.Request{Method: method, URL: roomURL, Body: `{"name":""}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Nil name value", Request: utils.Request{Method: method, URL: roomURL, Body: `{"name":nil}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Integer name value", Request: utils.Request{Method: method, URL: roomURL, Body: `{"name":1}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Object name value", Request: utils.Request{Method: method, URL: roomURL, Body: `{"name":{"somekey":"somevalue"}}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
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
			{Name: "Success", Request: utils.Request{Method: method, URL: roomURL}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"room":{"name":"test"}}`},
			{Name: "Not found", Request: utils.Request{Method: method, URL: roomURL + "0"}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
		},
	}
	test.Run(t)
}

// TestUpdateRoom test the UpdateRoom function.
func TestUpdateRoom(t *testing.T) {
	method := http.MethodPatch

	test := utils.TestRUD{
		CreateRequest:  roomCreateRequest,
		CreateResponse: CreateResponseRoom{},
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method, URL: roomURL, Body: `{"name":"test2"}`}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"room":{.+}}`},
			{Name: "Correctly updated", Request: utils.Request{Method: http.MethodGet, URL: roomURL}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"room":{"name":"test2"}}`},
			{Name: "Empty Body", Request: utils.Request{Method: method, URL: roomURL, Body: ``}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Empty json", Request: utils.Request{Method: method, URL: roomURL, Body: `{}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Bad name key", Request: utils.Request{Method: method, URL: roomURL, Body: `{"wrongkey":"test2"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Empty name value", Request: utils.Request{Method: method, URL: roomURL, Body: `{"name":""}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Nil name value", Request: utils.Request{Method: method, URL: roomURL, Body: `{"name":nil}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Integer name value", Request: utils.Request{Method: method, URL: roomURL, Body: `{"name":1}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Object name value", Request: utils.Request{Method: method, URL: roomURL, Body: `{"name":{"somekey":"somevalue"}}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Not found", Request: utils.Request{Method: method, URL: roomURL + "0", Body: `{"name":"test2"}`}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
		},
	}
	test.Run(t)
}

// TestDeleteRoom test the DeleteRoom function.
func TestDeleteRoom(t *testing.T) {
	method := http.MethodDelete

	test := utils.TestRUD{
		CreateRequest:  roomCreateRequest,
		CreateResponse: CreateResponseRoom{},
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method, URL: roomURL}, ResponseCode: http.StatusOK, ResponseBodyRegex: ``},
			{Name: "Correctly deleted", Request: utils.Request{Method: http.MethodGet, URL: roomURL}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
			{Name: "Not found", Request: utils.Request{Method: method, URL: roomURL + "0"}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
		},
	}
	test.Run(t)
}
