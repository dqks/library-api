package integration_test

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/handler"
	"library-api/internal/model"
	"library-api/internal/repository"
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
		// PUT request
		wantPutBody   string
		wantPutCode   int
		wantPutHeader http.Header
		// GET request
		wantGetBody string
	}{
		{
			name:        "success",
			id:          "1",
			body:        `{"title":"New edited title","author":"New edited author","year":2000,"available":false}`,
			wantPutBody: `{"id":1,"title":"New edited title","author":"New edited author","year":2000,"available":false}`,
			wantPutCode: 200,
			wantPutHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"New edited title","author":"New edited author","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:        "success, partial fields",
			id:          "1",
			body:        `{"title":"New edited title"}`,
			wantPutBody: `{"id":1,"title":"New edited title","author":"George Orwell","year":2000,"available":false}`,
			wantPutCode: 200,
			wantPutHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"New edited title","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:        "not found",
			id:          "5",
			body:        `{"title":"New edited title","author":"New edited author","year":2000,"available":false}`,
			wantPutBody: `{"error":"` + apperrors.ErrNotFound.Error() + `"}`,
			wantPutCode: 404,
			wantPutHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:        "zero id",
			id:          "0",
			body:        `{"title":"New edited title","author":"New edited author","year":2000,"available":false}`,
			wantPutBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantPutCode: 400,
			wantPutHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:        "negative id",
			id:          "-1",
			body:        `{"title":"New edited title","author":"New edited author","year":2000,"available":false}`,
			wantPutBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantPutCode: 400,
			wantPutHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:        "not number id",
			id:          "1hofg",
			body:        `{"title":"New edited title","author":"New edited author","year":2000,"available":false}`,
			wantPutBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantPutCode: 400,
			wantPutHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:          "cancelled",
			id:            "1",
			body:          ``,
			cancel:        true,
			wantPutBody:   ``,
			wantPutCode:   200,
			wantPutHeader: map[string][]string{},
			wantGetBody:   `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:            "deadline exceeded",
			id:              "1",
			body:            ``,
			deadlineContext: true,
			wantPutBody:     ``,
			wantPutCode:     504,
			wantPutHeader:   map[string][]string{},
			wantGetBody:     `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:        "bad json",
			id:          "1",
			body:        `{"title":"New edited title","author":"New edited author","year":2000,"available":`,
			wantPutBody: `{"error":"unexpected EOF"}`,
			wantPutCode: 400,
			wantPutHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
		{
			name:        "additional field",
			id:          "1",
			body:        `{"title":"New edited title","author":"New edited author","year":2000,"available":false,"additional":false}`,
			wantPutBody: `{"error":"json: unknown field \"additional\""}`,
			wantPutCode: 400,
			wantPutHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"1984","author":"George Orwell","year":2000,"available":false},{"id":2,"title":"New title","author":"New author","year":2000,"available":false}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repository.Create(
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
			mux.HandleFunc("/books/"+tt.id, handler.EditBookByID(service))
			mux.HandleFunc("/books", handler.GetBooks(service))

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
				http.MethodDelete,
				"/books/"+tt.id,
				strings.NewReader(tt.body),
			)
			req.SetPathValue("id", tt.id)

			// PUT request
			mux.ServeHTTP(rr, req)

			if tt.wantPutCode != rr.Code {
				t.Errorf(
					"expected code to be %d but got %d",
					tt.wantPutCode,
					rr.Code,
				)
			}

			if !testutil.CompareHeaders(tt.wantPutHeader, rr.Header()) {
				t.Errorf(
					"expected headers to be %v\nbut got %v",
					tt.wantPutHeader,
					rr.Header(),
				)
			}

			if strings.TrimSpace(tt.wantPutBody) != strings.TrimSpace(rr.Body.String()) {
				t.Errorf(
					"expected PUT body to be %v\nbut got %v",
					tt.wantPutBody,
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
