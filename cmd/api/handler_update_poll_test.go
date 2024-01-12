package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

func Test_app_updatePollHandler(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		json           string
		expectedStatus int
		expectedBody   string
	}{
		{"valid update partial", "1", `{"question":"changed"}`, http.StatusOK, `"question":"changed"`},
		{"valid update complete", "1", fmt.Sprintf(
			`{"question":"changed", "description":"added description", "expires_at":%q}`,
			time.Now().Add(2*time.Minute).Format(time.RFC3339),
		), http.StatusOK, `"question":"changed","description":"added description"`},
		{"empty json", "1", `{}`, http.StatusBadRequest, "no fields provided for update"},
		{"bad id", "-1", `{"question":"changed"}`, http.StatusBadRequest, "invalid id"},
		{"no record", "2", `{"question":"changed"}`, http.StatusNotFound, "the requested resource could not be found"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPatch, "/", strings.NewReader(test.json))
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("pollID", test.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.updatePollHandler)
			handler.ServeHTTP(rr, req)
			if rr.Code != test.expectedStatus {
				t.Errorf("expected status %d, but got %d", test.expectedStatus, rr.Code)
			}
			if !strings.Contains(rr.Body.String(), test.expectedBody) {
				t.Errorf("expected body to contain %q, but got %q", test.expectedBody, rr.Body)
			}
		})
	}
}
