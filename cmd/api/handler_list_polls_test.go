package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_app_listPollsHandler(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
		search         string
		sort           string
		page           any
		pageSize       int
		expectedBody   string
	}{
		{
			name:           "get polls default",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid page",
			expectedStatus: http.StatusUnprocessableEntity,
			page:           0,
			pageSize:       20,
			expectedBody:   `"page":"must be greater than zero"`,
		},
		{
			name:           "invalid page not int",
			expectedStatus: http.StatusUnprocessableEntity,
			page:           "a",
			pageSize:       20,
			expectedBody:   `"page":"must be an integer value"`,
		},
		{
			name:           "invalid page size",
			expectedStatus: http.StatusUnprocessableEntity,
			page:           1,
			pageSize:       0,
			expectedBody:   `"page_size":"must be greater than zero"`,
		},
		{
			name:           "invalid page size, over 50",
			expectedStatus: http.StatusUnprocessableEntity,
			page:           1,
			pageSize:       51,
			expectedBody:   `"page_size":"must be a maximum of 50"`,
		},
		{
			name:           "invalid sort value",
			expectedStatus: http.StatusUnprocessableEntity,
			sort:           "test",
			page:           1,
			pageSize:       20,
			expectedBody:   `"sort":"invalid sort value"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var url string

			if test.name == "get polls default" {
				url = "/polls"
			} else {
				url = fmt.Sprintf(
					"/polls?page=%v&page_size=%v&sort=%v&search=%v",
					test.page,
					test.pageSize,
					test.sort,
					test.search,
				)
			}
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.listPollsHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != test.expectedStatus {
				t.Errorf("expected status code %d, but got %d", test.expectedStatus, rr.Code)
			}

			if !strings.Contains(rr.Body.String(), test.expectedBody) {
				t.Errorf("expected body to contain %q, but got %q", test.expectedBody, rr.Body)
			}
		})
	}
}
