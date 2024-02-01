package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_app_rateLimit(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	app.config.limiter.rps = 2
	app.config.limiter.burst = 4
	app.config.limiter.enabled = true

	handlerToTest := app.rateLimit(nextHandler)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "0.0.0.0:0000"
	rr := httptest.NewRecorder()
	for i := 0; i < 6; i++ {
		handlerToTest.ServeHTTP(rr, req)
		if i < 4 && rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, but got %d", http.StatusOK, rr.Code)
		}
		if i > 4 && rr.Code != http.StatusTooManyRequests {
			t.Errorf("expected status code %d, but got %d", http.StatusTooManyRequests, rr.Code)
		}
	}

	app.config.limiter.enabled = false
	rr = httptest.NewRecorder()
	for i := 0; i < 10; i++ {
		handlerToTest.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, but got %d", http.StatusOK, rr.Code)
		}
	}
}
