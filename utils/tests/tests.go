//nolint:deadcode,unused
package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
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
	"gorm.io/gorm"
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
	// Name is the name of the test.
	Name string
	// Request is the request to send.
	Request Request
	// ResponseCode is the expected response code.
	ResponseCode int
	// ResponseBodyRegex is the expected response body regex.
	ResponseBodyRegex interface{}
}

// Request is a simple request to be sent to the router.
type Request struct {
	// Method is the HTTP method to use.
	Method string
	// Target is appended to the Test Target URL.
	Target string
	// Body is the body of the request.
	Body string
	// Headers is a list of headers to send with the request.
	// If not nil, they take precedence over the Test headers.
	// If nil, the Test headers are used instead.
	Headers []Header
}

// Header symbolizes a header to be sent with a request.
type Header struct {
	Key   string
	Value string
}

// CreateTestUser creates a new user for tests and returns the authorization header associated to the user.
func CreateTestUser(user models.User) (string, []Header, error) {
	var c utils.CreateResponse

	body, err := json.Marshal(user)
	if err != nil {
		return "", nil, err
	}
	res, err := executeRequest(http.MethodPost, "/users", string(body), nil)
	if err != nil {
		return "", nil, err
	}

	err = json.Unmarshal(res.Body.Bytes(), &c)
	if err != nil {
		return "", nil, err
	}
	id := path.Base(c.URI)

	return id, GetBasicAuthHeader(id, c.Password), nil
}

// GetBasicAuthHeader returns the Authorization header for a given id and password.
func GetBasicAuthHeader(id, password string) []Header {
	authHeader := base64.StdEncoding.EncodeToString([]byte(id + ":" + password))

	return []Header{
		{Key: "Authorization", Value: "Basic " + authHeader}}
}

// ResetTable resets a table in the database.
func ResetTable(db *gorm.DB, table ...interface{}) error {
	// Drop all tables.
	for _, t := range table {
		if db.Migrator().HasTable(t) {
			err := db.Migrator().DropTable(t)
			if err != nil {
				return err
			}
			err = db.AutoMigrate(t)
			if err != nil {
				return err
			}
		} else {
			return errors.New("table not found")
		}
	}
	return nil
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
