package integration_test

import (
	"context"
	"library-api/internal/handler"
	"library-api/internal/model"
	inmemoryrepo "library-api/internal/repository/inmemoryrepo"
	"library-api/internal/service"
	"library-api/internal/testutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetBooks(t *testing.T) {
	tests := []struct {
		name             string
		params           string
		cancel           bool
		deadlineExceeded bool
		// GET Book Request
		wantBody    string
		wantHeaders http.Header
		wantCode    int
	}{
		{
			name:   "success",
			params: "",
			// GET Book Request
			wantBody: `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":true},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 200,
		},
		{
			name:   "success only true",
			params: "?available=true",
			// GET Book Request
			wantBody: `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":true}]`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 200,
		},
		{
			name:   "success only false",
			params: "?available=false",
			// GET Book Request
			wantBody: `[{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 200,
		},
		{
			name:   "bad param",
			params: "?available=dasdaws",
			// GET Book Request
			wantBody: `{"error":"strconv.ParseBool: parsing \"dasdaws\": invalid syntax"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 400,
		},
		{
			name:   "cancel",
			params: "",
			cancel: true,
			// GET Book Request
			wantBody:    ``,
			wantHeaders: map[string][]string{},
			wantCode:    200,
		},
		{
			name:             "deadline exceeded",
			deadlineExceeded: true,
			params:           "",
			// GET Book Request
			wantBody:    ``,
			wantHeaders: map[string][]string{},
			wantCode:    504,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx context.Context
			var cancel context.CancelFunc

			if tt.deadlineExceeded {
				ctx, cancel = context.WithDeadline(context.Background(), time.Time{})
				defer cancel()
			} else {
				ctx, cancel = context.WithCancel(context.Background())
				defer cancel()
			}

			if tt.cancel {
				cancel()
			}

			repo := inmemoryrepo.Create(3,
				[]model.Book{
					{
						ID:        1,
						Title:     "1984",
						Author:    "George Orwell",
						Year:      2000,
						Available: true,
					},
					{
						ID:        2,
						Title:     "New title",
						Author:    "New author",
						Year:      2000,
						Available: false,
					},
				},
			)
			service := service.Create(repo)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /books", handler.GetBooks(service))

			rr := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(
				ctx,
				http.MethodGet,
				"/books"+tt.params,
				nil,
			)
			req.Header.Set("Content-Type", "application/json")

			mux.ServeHTTP(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf(
					"expected code to be %d but got %d",
					tt.wantCode,
					rr.Code,
				)
			}

			if !testutil.CompareHeaders(tt.wantHeaders, rr.Header()) {
				t.Errorf(
					"expected headers to be %v but got %v",
					tt.wantHeaders,
					rr.Header(),
				)
			}

			if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(tt.wantBody) {
				t.Errorf(
					"expected body to be %v\nbut got %v",
					tt.wantBody,
					rr.Body.String(),
				)
			}
		})
	}
}
