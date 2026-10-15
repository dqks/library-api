package tests

import (
	"context"
	"encoding/json"
	"errors"
	"library-api/internal/apperrors"
	"library-api/internal/handler"
	"library-api/internal/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetBookByID(t *testing.T) {
	tests := []struct {
		name                 string
		id                   string
		cancel               bool
		deadlineContext      bool
		serviceErr           error
		serviceBook          model.Book
		wantId               int
		wantBody             string
		wantHeaders          http.Header
		wantCode             int
		wantServiceCallCount int
	}{
		{
			name:       "success",
			id:         "1",
			serviceErr: nil,
			serviceBook: model.Book{
				ID:        1,
				Title:     "title",
				Author:    "author",
				Year:      2000,
				Available: false,
			},
			wantId:   1,
			wantBody: `{"id":1,"title":"title","author":"author","year":2000,"available":false}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode:             200,
			wantServiceCallCount: 1,
		},
		{
			name:        "not found",
			id:          "1",
			serviceErr:  apperrors.ErrNotFound,
			serviceBook: model.Book{},
			wantId:      1,
			wantBody:    `{"error":"` + apperrors.ErrNotFound.Error() + `"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode:             404,
			wantServiceCallCount: 1,
		},
		{
			name:        "deadline exceeded in service",
			id:          "1",
			serviceErr:  context.DeadlineExceeded,
			serviceBook: model.Book{},
			wantId:      1,
			wantBody:    `{"error":"` + context.DeadlineExceeded.Error() + `"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode:             504,
			wantServiceCallCount: 1,
		},
		{
			name:                 "canceled in service",
			id:                   "1",
			serviceErr:           context.Canceled,
			serviceBook:          model.Book{},
			wantId:               1,
			wantBody:             ``,
			wantHeaders:          map[string][]string{},
			wantCode:             200,
			wantServiceCallCount: 1,
		},
		{
			name:                 "canceled in handler",
			id:                   "1",
			cancel:               true,
			serviceErr:           nil,
			serviceBook:          model.Book{},
			wantId:               1,
			wantBody:             ``,
			wantHeaders:          map[string][]string{},
			wantCode:             200,
			wantServiceCallCount: 0,
		},
		{
			name:                 "canceled in handler",
			id:                   "1",
			deadlineContext:      true,
			serviceErr:           nil,
			serviceBook:          model.Book{},
			wantId:               1,
			wantBody:             ``,
			wantHeaders:          map[string][]string{},
			wantCode:             504,
			wantServiceCallCount: 0,
		},
		{
			name:        "id is zero",
			id:          "0",
			serviceErr:  nil,
			serviceBook: model.Book{},
			wantBody:    `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode:             400,
			wantServiceCallCount: 0,
		},
		{
			name:        "id is negative",
			id:          "-1",
			serviceErr:  nil,
			serviceBook: model.Book{},
			wantBody:    `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode:             400,
			wantServiceCallCount: 0,
		},
		{
			name:        "unknown service err",
			id:          "1",
			serviceErr:  errors.New("unknown"),
			serviceBook: model.Book{},
			wantId:      1,
			wantBody:    `{"error":"unknown"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode:             500,
			wantServiceCallCount: 1,
		},
		{
			name:        "not number id",
			id:          "1hgdkfrd",
			serviceErr:  errors.New("unknown"),
			serviceBook: model.Book{},
			wantBody:    `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode:             400,
			wantServiceCallCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			rr := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(
				ctx,
				http.MethodGet,
				"http://localhost:8080/books/"+tt.id,
				nil,
			)
			req.SetPathValue("id", tt.id)

			s := ServiceMock{
				err:  tt.serviceErr,
				book: tt.serviceBook,
			}

			handler.GetBookByID(&s)(rr, req)

			if tt.wantCode != rr.Code {
				t.Errorf(
					"expected code to be %d but got %d",
					tt.wantCode,
					rr.Code,
				)
			}

			if !compareHeaders(tt.wantHeaders, rr.Header()) {
				t.Errorf(
					"expected headers to be %v but got %v",
					tt.wantHeaders,
					rr.Header(),
				)
			}

			if s.callCount != tt.wantServiceCallCount {
				t.Errorf(
					"expected service to be called %d times but got %d",
					tt.wantServiceCallCount,
					s.callCount,
				)
			}

			if ctx != s.ctx && tt.wantServiceCallCount >= 1 {
				t.Error("expected context passed to service be equal to handlers's context")
			}

			if tt.wantId != s.id && tt.wantServiceCallCount >= 1 {
				t.Errorf(
					"expected id to be %d but got %d",
					tt.wantId,
					s.id,
				)
			}

			// Ожидаем вернуть в ответе книгу которую вернул сервис
			if tt.serviceErr == nil {
				var responseBook model.Book
				decoder := json.NewDecoder(strings.NewReader(tt.wantBody))
				err := decoder.Decode(&responseBook)
				if err != nil {
					t.Fatal("failed to decode body")
				}
				if responseBook != tt.serviceBook {
					t.Errorf(
						"expected book from repository to be the same in the response\nwanted book %v\n actual book %v",
						tt.serviceBook,
						responseBook,
					)
				}
			}

			if strings.TrimSpace(tt.wantBody) != strings.TrimSpace(rr.Body.String()) {
				t.Errorf(
					"expected body to be %v\nbut got %v",
					tt.wantBody,
					rr.Body.String(),
				)
			}
		})
	}
}
