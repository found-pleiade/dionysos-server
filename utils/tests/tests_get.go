package utils

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/assert/v2"
)

type ResponseCreate struct {
	ID string `json:"id"`
}

// TestGet struct to factorize code inside a single feature test.
type TestGet struct {
	CreateRequest Request
	SubTests      []SubTest
}

// Runs a series of tests for a Get type endpoint.
func (test TestGet) Run(t *testing.T) {
	disableLogs()

	// First create the ressource and get its ID.
	var created ResponseCreate
	w, err := executeRequest(test.CreateRequest.Method, test.CreateRequest.Url, test.CreateRequest.Body)
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(w.Body.Bytes(), &created)
	if err != nil {
		t.Error(err)
	}

	// Then run the tests.
	for _, subtest := range test.SubTests {
		t.Run(subtest.Name, func(t *testing.T) {
			url := subtest.Request.Url + created.ID

			w, err := executeRequest(subtest.Request.Method, url, subtest.Request.Body)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, w.Code, subtest.ResponseCode)
			assert.MatchRegex(t, w.Body.String(), subtest.ResponseBodyRegex)
		})
	}
}
