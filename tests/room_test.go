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

var roomUrl = "/rooms/"
var roomCreateRequest = utils.Request{Method: http.MethodPost, Url: roomUrl, Body: `{"name":"test"}`}

// TestCreateRoom tests the CreateRoom function.
func TestCreateRoom(t *testing.T) {
	method := http.MethodPost

	test := utils.TestCreate{
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method, Url: roomUrl, Body: `{"name":"test"}`}, ResponseCode: http.StatusCreated, ResponseBodyRegex: `{"id":"\d+"}`},
			{Name: "Empty body", Request: utils.Request{Method: method, Url: roomUrl, Body: ``}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Bad name key", Request: utils.Request{Method: method, Url: roomUrl, Body: `{"wrongkey":"test"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Bad name value", Request: utils.Request{Method: method, Url: roomUrl, Body: `{"name":""}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
		},
	}
	test.Run(t)
}

// TestGetRoom test the GetRoom function.
func TestGetRoom(t *testing.T) {
	method := http.MethodGet

	test := utils.TestRUD{
		CreateRequest:  roomCreateRequest,
		CreateResponse: CreateResponseRoom{},
		SubTests: []utils.SubTest{
			{Name: "Success", Request: utils.Request{Method: method, Url: roomUrl}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"room":{.+}}`},
			{Name: "Not found", Request: utils.Request{Method: method, Url: roomUrl + "0"}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
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
			{Name: "Success", Request: utils.Request{Method: method, Url: roomUrl, Body: `{"name":"test2"}`}, ResponseCode: http.StatusOK, ResponseBodyRegex: `{"room":{.+}}`},
			{Name: "Empty Body", Request: utils.Request{Method: method, Url: roomUrl, Body: ``}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Bad name key", Request: utils.Request{Method: method, Url: roomUrl, Body: `{"wrongkey":"test2"}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Bad name value", Request: utils.Request{Method: method, Url: roomUrl, Body: `{"name":""}`}, ResponseCode: http.StatusBadRequest, ResponseBodyRegex: `{"error":".+"}`},
			{Name: "Not found", Request: utils.Request{Method: method, Url: roomUrl + "0", Body: `{"name":"test2"}`}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
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
			{Name: "Success", Request: utils.Request{Method: method, Url: roomUrl}, ResponseCode: http.StatusOK, ResponseBodyRegex: ``},
			{Name: "Not found", Request: utils.Request{Method: method, Url: roomUrl + "0"}, ResponseCode: http.StatusNotFound, ResponseBodyRegex: `{"error":"Room not found"}`},
		},
	}
	test.Run(t)
}
