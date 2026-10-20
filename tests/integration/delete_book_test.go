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

func TestDeleteBook(t *testing.T) {

	tests := []struct {
		name            string
		id              string
		cancel          bool
		deadlineContext bool
		// DELETE Request
		wantDelBody    string
		wantDelCode    int
		wantDelHeaders http.Header
		// GET Request
		wantGetBody string
	}{
		{
			name:        "success",
			id:          "1",
			wantDelBody: `{"success":true}`,
			wantDelCode: 200,
			wantDelHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:        "not found",
			id:          "5",
			wantDelBody: `{"error":"` + apperrors.ErrNotFound.Error() + `"}`,
			wantDelCode: 404,
			wantDelHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:        "not number",
			id:          "1gedfl",
			wantDelBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantDelCode: 400,
			wantDelHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:        "zero id",
			id:          "0",
			wantDelBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantDelCode: 400,
			wantDelHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:        "negative id",
			id:          "-1",
			wantDelBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantDelCode: 400,
			wantDelHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:           "cancelled",
			id:             "1",
			cancel:         true,
			wantDelBody:    ``,
			wantDelCode:    200,
			wantDelHeaders: map[string][]string{},
			wantGetBody:    `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:            "deadline exceeded",
			id:              "1",
			deadlineContext: true,
			wantDelBody:     ``,
			wantDelCode:     504,
			wantDelHeaders:  map[string][]string{},
			wantGetBody:     `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := inmemory.Create(
				3,
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
			mux.HandleFunc("DELETE /books/{id}", handler.DeleteBookByID(service))
			mux.HandleFunc("GET /books", handler.GetBooks(service))

			// DELETE Request

			var ctx context.Context
			var cancel context.CancelFunc

			if tt.deadlineContext {
				ctx, cancel = context.WithDeadline(context.Background(), time.Time{})
				defer cancel()
			} else {
				ctx, cancel = context.WithCancel(context.Background())
				defer cancel()
			}

			if tt.cancel {
				cancel()
			}

			rrDel := httptest.NewRecorder()
			reqDel := httptest.NewRequestWithContext(
				ctx, http.MethodDelete,
				"/books/"+tt.id,
				nil,
			)

			mux.ServeHTTP(rrDel, reqDel)

			if tt.wantDelCode != rrDel.Code {
				t.Errorf(
					"expected code to be %d but got %d",
					tt.wantDelCode,
					rrDel.Code,
				)
			}

			if !testutil.CompareHeaders(tt.wantDelHeaders, rrDel.Header()) {
				t.Errorf(
					"expected headers to be %v\nbut got %v",
					tt.wantDelHeaders,
					rrDel.Header(),
				)
			}

			if strings.TrimSpace(tt.wantDelBody) != strings.TrimSpace(rrDel.Body.String()) {
				t.Errorf(
					"expected body to be %v\nbut got %v",
					tt.wantDelBody,
					rrDel.Body.String(),
				)
			}

			// GET Request

			ctx = context.WithoutCancel(context.Background())
			rrGet := httptest.NewRecorder()
			reqGet := httptest.NewRequestWithContext(
				ctx, http.MethodGet,
				"/books",
				nil,
			)

			mux.ServeHTTP(rrGet, reqGet)

			if strings.TrimSpace(tt.wantGetBody) != strings.TrimSpace(rrGet.Body.String()) {
				t.Errorf(
					"expected body to be %v\nbut got %v",
					tt.wantDelBody,
					rrGet.Body.String(),
				)
			}
		})
	}
}
