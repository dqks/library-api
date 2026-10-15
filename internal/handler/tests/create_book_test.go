package tests

import (
	"context"
	"encoding/json"
	"errors"
	"library-api/internal/apperrors"
	"library-api/internal/handler"
	"library-api/internal/model"
	"library-api/internal/service"
	"library-api/internal/testutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestCreateBook(t *testing.T) {
	tests := []struct {
		name             string
		body             string
		serviceCallCount int
		serviceErr       error
		serviceBook      model.Book
		deadlineContext  bool
		cancel           bool
		wantCode         int
		wantHeader       http.Header
		wantBody         string
	}{
		{
			name:             "success",
			body:             `{"title":"New title","author":"New Author","year":2000,"available":false}`,
			serviceCallCount: 1,
			wantCode:         201,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			serviceErr: nil,
			serviceBook: model.Book{
				ID:        1,
				Title:     "New title",
				Author:    "New Author",
				Year:      2003,
				Available: false,
			},
			wantBody: `{"id":1,"title":"New title","author":"New Author","year":2003,"available":false}`,
		},
		{
			name:             "canceled",
			body:             `{"title":"New title","author":"New Author","year":2000,"available":false}`,
			serviceCallCount: 0,
			wantCode:         200,
			wantHeader:       map[string][]string{},
			cancel:           true,
			serviceErr:       nil,
			serviceBook:      model.Book{},
			wantBody:         ``,
		},
		{
			name:             "err, body invalid values",
			body:             `{"title":"New title","author":"New Author","year":2077,"available":false}`,
			serviceCallCount: 1,
			wantCode:         400,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			serviceErr:  apperrors.ErrInvalidValues,
			serviceBook: model.Book{},
			wantBody:    `{"error":"` + apperrors.ErrInvalidValues.Error() + `"}`,
		},
		{
			name:             "err, body absence of required field",
			body:             `{"title":"New title","author":"New Author","available":false}`,
			serviceCallCount: 1,
			wantCode:         400,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			serviceErr:  apperrors.ErrRequiredFields,
			serviceBook: model.Book{},
			wantBody:    `{"error":"` + apperrors.ErrRequiredFields.Error() + `"}`,
		},
		{
			name:             "err, incorrect json",
			body:             `{"title":"New title","author":"New Author","available":`,
			serviceCallCount: 0,
			wantCode:         400,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			serviceErr:  nil,
			serviceBook: model.Book{},
			wantBody:    `{"error":"unexpected EOF"}`,
		},
		{
			name:             "err, unknown field json",
			body:             `{"title":"New title","author":"New Author","available":false,"test":"test"}`,
			serviceCallCount: 0,
			wantCode:         400,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			serviceErr:  nil,
			serviceBook: model.Book{},
			wantBody:    `{"error":"json: unknown field \"test\""}`,
		},
		{
			name:             "err, deadline exceeded",
			body:             `{"title":"New title","author":"New Author","available":false}`,
			serviceCallCount: 0,
			wantCode:         504,
			deadlineContext:  true,
			wantHeader:       map[string][]string{},
			serviceErr:       nil,
			serviceBook:      model.Book{},
			wantBody:         ``,
		},
		{
			name:             "err, service canceled context",
			body:             `{"title":"New title","author":"New Author","year":2000,"available":false}`,
			serviceCallCount: 1,
			wantCode:         200,
			wantHeader:       map[string][]string{},
			serviceErr:       context.Canceled,
			serviceBook:      model.Book{},
			wantBody:         ``,
		},
		{
			name:             "err, deadline exceeded in service",
			body:             `{"title":"New title","author":"New Author","year":2000,"available":false}`,
			serviceCallCount: 1,
			wantCode:         504,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			serviceErr:  context.DeadlineExceeded,
			serviceBook: model.Book{},
			wantBody:    `{"error":"` + context.DeadlineExceeded.Error() + `"}`,
		},
		{
			name:             "internal err",
			body:             `{"title":"New title","author":"New Author","year":2000,"available":false}`,
			serviceCallCount: 1,
			wantCode:         500,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			serviceErr:  apperrors.ErrInternal,
			serviceBook: model.Book{},
			wantBody:    `{"error":"` + apperrors.ErrInternal.Error() + `"}`,
		},
		{
			name:             "unknown err",
			body:             `{"title":"New title","author":"New Author","year":2000,"available":false}`,
			serviceCallCount: 1,
			wantCode:         500,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			serviceErr:  errors.New("new error"),
			serviceBook: model.Book{},
			wantBody:    `{"error":"new error"}`,
		},
		{
			name:             "err, wrong year type",
			body:             `{"title":"New title","author":"New Author","available":false,"year":"2000"}`,
			serviceCallCount: 0,
			wantCode:         400,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			serviceErr:  nil,
			serviceBook: model.Book{},
			wantBody:    `{"error":"json: cannot unmarshal string into Go struct field CreateBookBody.year of type uint16"}`,
		},
		{
			name:             "err, wrong available type",
			body:             `{"title":"New title","author":"New Author","available":"false","year":2000}`,
			serviceCallCount: 0,
			wantCode:         400,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			serviceErr:  nil,
			serviceBook: model.Book{},
			wantBody:    `{"error":"json: cannot unmarshal string into Go struct field CreateBookBody.available of type bool"}`,
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

			req := httptest.NewRequestWithContext(
				ctx,
				http.MethodPost,
				"http://localhost:8080/books",
				strings.NewReader(tt.body),
			)

			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			s := ServiceMock{
				err:  tt.serviceErr,
				book: tt.serviceBook,
			}

			if tt.cancel {
				cancel()
			}

			handler.CreateBook(&s)(rr, req)

			if ctx != s.ctx && tt.serviceCallCount >= 1 {
				t.Fatalf("context passed to service is not equal to request's context")
			}

			if rr.Code != tt.wantCode {
				t.Fatalf("expected code %d but got %d", tt.wantCode, rr.Code)
			}

			if !testutil.CompareHeaders(tt.wantHeader, rr.Header()) {
				t.Fatalf(
					"expected headers to be %v\nbut got %v",
					tt.wantHeader,
					rr.Header(),
				)
			}

			if tt.serviceCallCount != s.callCount {
				t.Fatalf(
					"expected to call service %d times but got %d",
					tt.serviceCallCount,
					s.callCount,
				)
			}

			if s.callCount >= 1 {
				var expectedInput service.CreateBookInput
				decoder := json.NewDecoder(strings.NewReader(tt.body))
				err := decoder.Decode(&expectedInput)
				if err != nil {
					t.Fatal("failed to decode expected input")
				}
				if !reflect.DeepEqual(expectedInput, s.createInput) {
					t.Fatalf(
						"expected to call service with input %v\nbut got %v",
						expectedInput,
						s.createInput,
					)
				}
			}

			if strings.TrimSpace(tt.wantBody) != strings.TrimSpace(rr.Body.String()) {
				t.Fatalf(
					"expected result body %s\nbut got %s",
					tt.wantBody,
					rr.Body.String(),
				)
			}
		})
	}
}
