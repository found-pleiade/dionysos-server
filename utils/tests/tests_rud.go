//nolint:typecheck
package utils

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

// TestRUD (Read, Update, Delete) struct to factorize code inside a single feature test.
type TestRUD struct {
	// CreateRequest creates the object against which the tests will run.
	CreateRequest Request

	// CreateRequestHeaders is a list of headers to send with the Create request.
	CreateRequestHeaders []Header

	// CreateResponse parses the response from the CreateRequest.
	// The returned string is used as a prefix URL for SubTests Request.Target.
	CreateResponse ICreateResponse

	// Request.Target field is appended to the CreateResponse.UriCreated() to create the full URL.
	SubTests []SubTest
}

// Runs a series of tests for the Get/Update/Delete type endpoint.
func (test TestRUD) Run(t *testing.T) {
	disableLogs()

	w, err := executeRequest(test.CreateRequest.Method, test.CreateRequest.Target, test.CreateRequest.Body, test.CreateRequestHeaders)
	if err != nil {
		t.Error(err)
	}

	uri, err := test.CreateResponse.TargetURI(w.Body.Bytes())
	if err != nil {
		t.Error(err)
	}
	headers, err := test.CreateResponse.Headers(w.Body.Bytes())
	if err != nil {
		t.Error(err)
	}

	// Then run the tests.
	for _, subtest := range test.SubTests {
		t.Run(subtest.Name, func(t *testing.T) {
			url := uri + subtest.Request.Target

			w, err := executeRequest(subtest.Request.Method, url, subtest.Request.Body, append(headers, subtest.Headers...))
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, w.Code, subtest.ResponseCode)
			assert.MatchRegex(t, w.Body.String(), subtest.ResponseBodyRegex)
		})
	}
}
