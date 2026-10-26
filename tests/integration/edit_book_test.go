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

func TestEditBook(t *testing.T) {
	tests := []struct {
		name            string
		cancel          bool
		deadlineContext bool
		id              string
		body            string
		// PATCH request
		wantPatchBody   string
		wantPatchCode   int
		wantPatchHeader http.Header
		// GET request
		wantGetBody string
	}{
		{
			name:          "success",
			id:            "1",
			body:          `{"title":"New edited title","author":"New edited author","year":2000,"available":false}`,
			wantPatchBody: `{"id":1,"title":"New edited title","author":"New edited author","year":2000,"available":false}`,
			wantPatchCode: 200,
			wantPatchHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"New edited title","author":"New edited author","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:          "success, partial fields",
			id:            "1",
			body:          `{"title":"New edited title"}`,
			wantPatchBody: `{"id":1,"title":"New edited title","author":"George Orwell","year":2000,"available":false}`,
			wantPatchCode: 200,
			wantPatchHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"New edited title","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:          "not found",
			id:            "5",
			body:          `{"title":"New edited title","author":"New edited author","year":2000,"available":false}`,
			wantPatchBody: `{"error":"` + apperrors.ErrNotFound.Error() + `"}`,
			wantPatchCode: 404,
			wantPatchHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:          "zero id",
			id:            "0",
			body:          `{"title":"New edited title","author":"New edited author","year":2000,"available":false}`,
			wantPatchBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantPatchCode: 400,
			wantPatchHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:          "negative id",
			id:            "-1",
			body:          `{"title":"New edited title","author":"New edited author","year":2000,"available":false}`,
			wantPatchBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantPatchCode: 400,
			wantPatchHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:          "not number id",
			id:            "1hofg",
			body:          `{"title":"New edited title","author":"New edited author","year":2000,"available":false}`,
			wantPatchBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantPatchCode: 400,
			wantPatchHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:            "cancelled",
			id:              "1",
			body:            ``,
			cancel:          true,
			wantPatchBody:   ``,
			wantPatchCode:   200,
			wantPatchHeader: map[string][]string{},
			wantGetBody:     `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:            "deadline exceeded",
			id:              "1",
			body:            ``,
			deadlineContext: true,
			wantPatchBody:   ``,
			wantPatchCode:   504,
			wantPatchHeader: map[string][]string{},
			wantGetBody:     `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:          "bad json",
			id:            "1",
			body:          `{"title":"New edited title","author":"New edited author","year":2000,"available":`,
			wantPatchBody: `{"error":"unexpected EOF"}`,
			wantPatchCode: 400,
			wantPatchHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:          "additional field",
			id:            "1",
			body:          `{"title":"New edited title","author":"New edited author","year":2000,"available":false,"additional":false}`,
			wantPatchBody: `{"error":"json: unknown field \"additional\""}`,
			wantPatchCode: 400,
			wantPatchHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
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

			service := service.CreateBookService(repo)

			mux := http.NewServeMux()
			mux.HandleFunc("PATCH /books/{id}", handler.EditBookByID(service))
			mux.HandleFunc("GET /books", handler.GetBooks(service))

			var ctxGet context.Context
			var cancel context.CancelFunc

			if tt.deadlineContext {
				ctxGet, cancel = context.WithDeadline(context.Background(), time.Time{})
				defer cancel()
			} else {
				ctxGet, cancel = context.WithCancel(context.Background())
				defer cancel()
			}

			if tt.cancel {
				cancel()
			}

			rr := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(
				ctxGet,
				http.MethodPatch,
				"/books/"+tt.id,
				strings.NewReader(tt.body),
			)

			// PUT request
			mux.ServeHTTP(rr, req)

			if tt.wantPatchCode != rr.Code {
				t.Errorf(
					"expected code to be %d but got %d",
					tt.wantPatchCode,
					rr.Code,
				)
			}

			if !testutil.CompareHeaders(tt.wantPatchHeader, rr.Header()) {
				t.Errorf(
					"expected headers to be %v\nbut got %v",
					tt.wantPatchHeader,
					rr.Header(),
				)
			}

			if strings.TrimSpace(tt.wantPatchBody) != strings.TrimSpace(rr.Body.String()) {
				t.Errorf(
					"expected PUT body to be %v\nbut got %v",
					tt.wantPatchBody,
					rr.Body.String(),
				)
			}

			// GET request
			ctxGet = context.WithoutCancel(context.Background())
			rrGet := httptest.NewRecorder()
			reqGet := httptest.NewRequestWithContext(
				ctxGet,
				http.MethodGet,
				"/books",
				nil,
			)

			mux.ServeHTTP(rrGet, reqGet)

			if strings.TrimSpace(tt.wantGetBody) != strings.TrimSpace(rrGet.Body.String()) {
				t.Errorf(
					"expected GET body to be %v\nbut got %v",
					tt.wantGetBody,
					rrGet.Body.String(),
				)
			}
		})
	}
}
