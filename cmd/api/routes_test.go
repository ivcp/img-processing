package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_app_routes(t *testing.T) {
	tests := []struct {
		route  string
		method string
	}{
		{"/v1/healthcheck", http.MethodGet},
		{"/v1/polls", http.MethodPost},
		{"/v1/polls/{id}", http.MethodGet},
		{"/v1/polls/{id}", http.MethodPut},
	}
	testMux := app.routes()
	chiRoutes := testMux.(chi.Routes)
	for _, test := range tests {
		if !routeExists(test.route, test.method, chiRoutes) {
			t.Errorf("route %q is not registered", test.route)
		}
	}
}

func routeExists(testRoute, testMethod string, chiRoutes chi.Routes) bool {
	found := false
	_ = chi.Walk(chiRoutes, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.EqualFold(method, testMethod) && strings.EqualFold(route, testRoute) {
			found = true
		}
		return nil
	})
	return found
}
