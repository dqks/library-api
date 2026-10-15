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

func TestCreateBook(t *testing.T) {
	tests := []struct {
		// init data
		name                string
		postBody            string
		cancelPost          bool
		deadlineContextPost bool
		// WANT POST
		wantPostBody    string
		wantPostCode    int
		wantPostHeaders http.Header
		// WANT GET
		wantGetBody string
	}{
		{
			name:         "success",
			postBody:     `{"title":"new title","author":"new author","year":2000,"available":true}`,
			wantPostBody: `{"id":2,"title":"new title","author":"new author","year":2000,"available":true}`,
			wantPostCode: 201,
			wantPostHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"1984","author":"Gerorge Orwell","year":2000,"available":true},{"id":2,"title":"new title","author":"new author","year":2000,"available":true}]`,
		},
		{
			name:         "success with Title instead of title",
			postBody:     `{"Title":"new title","author":"new author","year":2000,"available":true}`,
			wantPostBody: `{"id":2,"title":"new title","author":"new author","year":2000,"available":true}`,
			wantPostCode: 201,
			wantPostHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"1984","author":"Gerorge Orwell","year":2000,"available":true},{"id":2,"title":"new title","author":"new author","year":2000,"available":true}]`,
		},
		{
			name:            "canceled",
			cancelPost:      true,
			postBody:        `{"title":"new title","author":"new author","year":2000,"available":true}`,
			wantPostBody:    ``,
			wantPostCode:    200,
			wantPostHeaders: map[string][]string{},
			wantGetBody:     `[{"id":1,"title":"1984","author":"Gerorge Orwell","year":2000,"available":true}]`,
		},
		{
			name:                "deadline exceeded",
			deadlineContextPost: true,
			postBody:            `{"title":"new title","author":"new author","year":2000,"available":true}`,
			wantPostBody:        ``,
			wantPostCode:        504,
			wantPostHeaders:     map[string][]string{},
			wantGetBody:         `[{"id":1,"title":"1984","author":"Gerorge Orwell","year":2000,"available":true}]`,
		},
		{
			name:         "missing required field",
			postBody:     `{"title":"new title","year":2000,"available":true}`,
			wantPostBody: `{"error":"` + apperrors.ErrRequiredFields.Error() + `"}`,
			wantPostCode: 400,
			wantPostHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"1984","author":"Gerorge Orwell","year":2000,"available":true}]`,
		},
		{
			name:         "additional field",
			postBody:     `{"title":"new title","author":"new author","year":2000,"available":true,"additional":123}`,
			wantPostBody: `{"error":"json: unknown field \"additional\""}`,
			wantPostCode: 400,
			wantPostHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBody: `[{"id":1,"title":"1984","author":"Gerorge Orwell","year":2000,"available":true}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repository.Create(
				2, []model.Book{
					{
						ID:        1,
						Title:     "1984",
						Author:    "Gerorge Orwell",
						Year:      2000,
						Available: true,
					},
				})

			service := service.Create(repo)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /books", handler.GetBooks(service))
			mux.HandleFunc("POST /books", handler.CreateBook(service))

			// POST request
			var ctxPost context.Context
			var cancelPost context.CancelFunc

			if tt.deadlineContextPost {
				ctxPost, cancelPost = context.WithDeadline(context.Background(), time.Time{})
				defer cancelPost()
			} else {
				ctxPost, cancelPost = context.WithCancel(context.Background())
				defer cancelPost()
			}

			if tt.cancelPost {
				cancelPost()
			}

			req := httptest.NewRequestWithContext(
				ctxPost,
				http.MethodPost,
				"/books",
				strings.NewReader(tt.postBody),
			)
			req.Header.Set("Content-Type", "application/json")

			rrPost := httptest.NewRecorder()
			mux.ServeHTTP(rrPost, req)

			if rrPost.Code != tt.wantPostCode {
				t.Errorf("expected code to be %d but got %d", tt.wantPostCode, rrPost.Code)
			}

			if !testutil.CompareHeaders(tt.wantPostHeaders, rrPost.Header()) {
				t.Errorf(
					"expected headers to be %v\nbut got %v",
					tt.wantPostHeaders,
					rrPost.Header(),
				)
			}

			if strings.TrimSpace(tt.wantPostBody) != strings.TrimSpace(rrPost.Body.String()) {
				t.Errorf(
					"expected POST body to be %v\nbut got %v",
					tt.wantPostBody,
					rrPost.Body.String(),
				)
			}

			// GET request
			rrGet := httptest.NewRecorder()
			ctxGet := context.WithoutCancel(context.Background())

			req = httptest.NewRequestWithContext(
				ctxGet,
				http.MethodGet,
				"/books",
				nil,
			)

			mux.ServeHTTP(rrGet, req)

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
