package integration_test

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/handler"
	"library-api/internal/model"
	"library-api/internal/repository/inmemory"
	"library-api/internal/service"
	"library-api/internal/testutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetBook(t *testing.T) {
	tests := []struct {
		name             string
		id               string
		cancel           bool
		deadlineExceeded bool
		wantBody         string
		wantHeaders      http.Header
		wantCode         int
	}{
		{
			name: "success",
			id:   "1",
			//GET Book Request
			wantBody: `{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 200,
		},
		{
			name: "not found",
			id:   "5",
			//GET Book Request
			wantBody: `{"error":"` + apperrors.ErrNotFound.Error() + `"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 404,
		},
		{
			name: "zero id",
			id:   "0",
			//GET Book Request
			wantBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 400,
		},
		{
			name: "negative id",
			id:   "-1",
			//GET Book Request
			wantBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 400,
		},
		{
			name: "id not number",
			id:   "1dlhfkg",
			//GET Book Request
			wantBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 400,
		},
		{
			name:   "cancelled",
			id:     "1",
			cancel: true,
			//GET Book Request
			wantBody:    ``,
			wantHeaders: map[string][]string{},
			wantCode:    200,
		},
		{
			name:             "deadlne exceeded",
			id:               "1",
			deadlineExceeded: true,
			//GET Book Request
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

			repo := inmemory.Create(3,
				[]model.Book{
					{
						ID:        1,
						Title:     "1984",
						Author:    "George Orwell",
						Year:      2000,
						Available: false,
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
			mux.HandleFunc("GET /books/{id}", handler.GetBookByID(service))

			rr := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(
				ctx,
				http.MethodGet,
				"/books/"+tt.id,
				nil,
			)

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
