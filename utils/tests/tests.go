//nolint:deadcode,unused
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
