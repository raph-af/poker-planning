package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// these tests are flaky

func TestLimitRate(t *testing.T) {
	totalCalls := 3
	responsesCodes := make(chan int, totalCalls)

	for i := 0; i < totalCalls; i++ {
		go func() {
			inputHandler := func(w http.ResponseWriter, r *http.Request) {}
			req := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
			res := httptest.NewRecorder()
			inputHandler(res, req)

			outputHandler := LimitRate(inputHandler)
			outputHandler.ServeHTTP(res, req)

			responsesCodes <- res.Code
		}()
	}

	hasStatusTooManyRequests := false
	for i := 0; i < totalCalls; i++ {
		if <-responsesCodes == http.StatusTooManyRequests {
			hasStatusTooManyRequests = true
		}
	}
	if !hasStatusTooManyRequests {
		t.Errorf(`Expected at least one response with %v code out of %v successive calls`,
			http.StatusTooManyRequests, totalCalls)
	}
}

func TestLimitRateByIp(t *testing.T) {
	totalCalls := 3
	responsesCodes := make(chan int, totalCalls)

	for i := 0; i < totalCalls; i++ {
		go func() {
			inputHandler := func(w http.ResponseWriter, r *http.Request) {}
			req := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
			req.RemoteAddr = "127.0.0.1"
			res := httptest.NewRecorder()

			inputHandler(res, req)

			outputHandler := LimitRateByIp(inputHandler)
			outputHandler.ServeHTTP(res, req)

			responsesCodes <- res.Code
		}()
	}

	hasStatusTooManyRequests := false
	for i := 0; i < totalCalls; i++ {
		if <-responsesCodes == http.StatusTooManyRequests {
			hasStatusTooManyRequests = true
		}
	}
	if !hasStatusTooManyRequests {
		t.Errorf(`Expected at least one response with %v code out of %v successive calls`,
			http.StatusTooManyRequests, totalCalls)
	}
}
