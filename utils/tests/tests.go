package utils

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Brawdunoir/dionysos-server/routes"
	"github.com/gin-gonic/gin"
)

var router = routes.SetupRouter(gin.New())

type Test interface {
	Run(t *testing.T)
}

// SubTest to be send to the router during series of tests.
type SubTest struct {
	Name              string
	Request           Request
	ResponseCode      int
	ResponseBodyRegex interface{}
}

type Request struct {
	Method string
	Url    string
	Body   string
}

// disableLogs to remove logs from default logger during tests.
func disableLogs() {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
}

// executeTest executes a single request and returns the response.
func executeRequest(method, url, body string) (w *httptest.ResponseRecorder, err error) {
	w = httptest.NewRecorder()
	req, err := http.NewRequest(method, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}

	router.ServeHTTP(w, req)
	return
}
