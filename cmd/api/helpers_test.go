package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_app_readIDParam(t *testing.T) {
	tests := []struct {
		name        string
		paramId     string
		expectError bool
	}{
		{"valid id", "1", false},
		{"invalid id", "a", true},
		{"invalid id", "-5", true},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", test.paramId)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		_, err := app.readIDParam(req)
		if !test.expectError && err != nil {
			t.Errorf("%s: expected no err, but got one: %q", test.name, err)
		}
		if test.expectError && err == nil {
			t.Errorf("%s: expected err, but didn't get one", test.name)
		}
	}
}
