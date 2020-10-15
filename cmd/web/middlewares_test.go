package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// these tests are flaky

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
		t.Errorf(`Expected at least one response with %v code out of %v successive calls`,
			http.StatusTooManyRequests, maxCalls)
	}
}

func TestLimitRateByIp(t *testing.T) {
	maxCalls := 3
	ip := "127.0.0.1"

	inputHandler := func(w http.ResponseWriter, r *http.Request) {}

	hasStatusTooManyRequests := false
	i := 0
	for i < maxCalls {
		req := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
		req.RemoteAddr = ip
		res := httptest.NewRecorder()
		inputHandler(res, req)

		outputHandler := LimitRateByIp(inputHandler)
		outputHandler.ServeHTTP(res, req)

		if res.Code == http.StatusTooManyRequests {
			hasStatusTooManyRequests = true
		}

		i++
	}

	if !hasStatusTooManyRequests {
		t.Errorf(`Expected at least one response with %v code out of %v successive calls`,
			http.StatusTooManyRequests, maxCalls)
	}
}
