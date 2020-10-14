package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// this test is flaky
func TestLimitRate(t *testing.T) {
	maxCalls := 3

	inputHandler := func(w http.ResponseWriter, r *http.Request) {}

	hasStatusTooManyRequests := false
	i := 0
	for i < maxCalls {
		req := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
		res := httptest.NewRecorder()
		inputHandler(res, req)

		outputHandler := LimitRate(inputHandler)
		outputHandler.ServeHTTP(res, req)

		if res.Code == http.StatusTooManyRequests {
			hasStatusTooManyRequests = true
		}

		i++
	}

	if !hasStatusTooManyRequests {
		t.Errorf(`Expected at least one %v response out of %v successive calls`, maxCalls, http.StatusTooManyRequests)
	}
}
