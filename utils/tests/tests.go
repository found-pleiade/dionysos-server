//nolint:deadcode,unused
package utils

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/Brawdunoir/dionysos-server/models"
	"github.com/Brawdunoir/dionysos-server/routes"
	utils "github.com/Brawdunoir/dionysos-server/utils/routes"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

var router = routes.SetupRouter(gin.New())

// ITest interface to run tests.
type ITest interface {
	Run(t *testing.T)
}

// ICreateResponse allows to map a Create request and retrieve the URI for further tests.
type ICreateResponse interface {
	TargetURI([]byte) (string, error)
	Headers([]byte) ([]Header, error)
}

// SubTest is an atomic test that includes a request and its intended response.
type SubTest struct {
	Name              string
	Request           Request
	Headers           []Header
	ResponseCode      int
	ResponseBodyRegex interface{}
}

// Request a simple request to be sent to the router.
type Request struct {
	Method string
	Target string
	Body   string
}

type Header struct {
	Key   string
	Value string
}

// CreateTestUser creates a new user for tests and returns the authorization header associated to the user.
func CreateTestUser(user models.User) ([]Header, error) {
	var c utils.CreateResponse

	body, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	res, err := executeRequest(http.MethodPost, "/users", string(body), nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res.Body.Bytes(), &c)
	if err != nil {
		return nil, err
	}

	id := path.Base(c.URI)
	password := c.Password
	authHeader := base64.StdEncoding.EncodeToString([]byte(id + ":" + password))

	return []Header{
		{Key: "Authorization", Value: "Basic " + authHeader}}, nil
}

// disableLogs to remove logs from default logger during tests.
func disableLogs() {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
}

// executeTest executes a single request and returns the response.
func executeRequest(method, url, body string, headers []Header) (w *httptest.ResponseRecorder, err error) {
	w = httptest.NewRecorder()
	req, err := http.NewRequest(method, url, bytes.NewBufferString(body))
	for _, header := range headers {
		req.Header.Set(header.Key, header.Value)
	}

	if err != nil {
		return nil, err
	}

	router.ServeHTTP(w, req)
	return
}
