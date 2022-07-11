//nolint:typecheck
package utils

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

// TestCreate struct to factorize code inside a single feature test.
type TestCreate struct {
	SubTests []SubTest
}

// Runs a series of tests for the Create type endpoint.
func (test TestCreate) Run(t *testing.T) {
	disableLogs()

	for _, subtest := range test.SubTests {
		t.Run(subtest.Name, func(t *testing.T) {
			w, err := executeRequest(subtest.Request.Method, subtest.Request.URL, subtest.Request.Body)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, w.Code, subtest.ResponseCode)
			assert.MatchRegex(t, w.Body.String(), subtest.ResponseBodyRegex)
		})
	}
}
