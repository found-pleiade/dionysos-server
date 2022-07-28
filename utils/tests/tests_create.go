//nolint:typecheck
package utils

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

// TestCreate struct to factorize code inside a single feature test.
type TestCreate struct {
	// Target is the URL of the endpoint to create the resource.
	Target string

	// Headers is a list of headers to send with all the requests.
	Headers []Header

	// SubTests is a list of SubTest to run.
	// If Request.Target is empty, the Target is used.
	// If Request.Target is not empty, the Request.Target is used instead.
	SubTests []SubTest
}

// Runs a series of tests for the Create type endpoint.
func (test TestCreate) Run(t *testing.T) {
	disableLogs()

	for _, subtest := range test.SubTests {
		var url string

		if subtest.Request.Target != "" {
			url = subtest.Request.Target
		} else {
			url = test.Target
		}

		t.Run(subtest.Name, func(t *testing.T) {
			w, err := executeRequest(subtest.Request.Method, url, subtest.Request.Body, append(subtest.Headers, test.Headers...))
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, w.Code, subtest.ResponseCode)
			assert.MatchRegex(t, w.Body.String(), subtest.ResponseBodyRegex)
		})
	}
}
