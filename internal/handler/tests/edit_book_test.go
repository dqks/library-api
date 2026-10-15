package tests

import (
	"context"
	"encoding/json"
	"library-api/internal/apperrors"
	"library-api/internal/handler"
	"library-api/internal/model"
	"library-api/internal/service"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestEditBookByID(t *testing.T) {
	tests := []struct {
		name                 string
		id                   string
		body                 string
		cancel               bool
		deadlineContext      bool
		serviceErr           error
		serviceBook          model.Book
		wantServiceCallCount int
		wantId               int
		wantErr              error
		wantBody             string
		wantHeaders          http.Header
		wantCode             int
	}{
		{
			name:       "success",
			id:         "1",
			body:       `{"title":"New title","author":"New Author","year":2000,"available":false}`,
			serviceErr: nil,
			serviceBook: model.Book{
				ID:        1,
				Title:     "New title",
				Author:    "New Author",
				Year:      2000,
				Available: false,
			},
			wantServiceCallCount: 1,
			wantId:               1,
			wantErr:              nil,
			wantBody:             `{"id":1,"title":"New title","author":"New Author","year":2000,"available":false}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 200,
		},
		{
			name:       "success, partial body",
			id:         "1",
			body:       `{"title":"New title"}`,
			serviceErr: nil,
			serviceBook: model.Book{
				ID:        1,
				Title:     "New title",
				Author:    "New Author",
				Year:      2000,
				Available: false,
			},
			wantServiceCallCount: 1,
			wantId:               1,
			wantErr:              nil,
			wantBody:             `{"id":1,"title":"New title","author":"New Author","year":2000,"available":false}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 200,
		},
		{
			name:                 "err, id is zero",
			id:                   "0",
			body:                 `{"title":"New title","author":"New Author","year":2000,"available":false}`,
			serviceErr:           nil,
			serviceBook:          model.Book{},
			wantServiceCallCount: 0,
			wantErr:              nil,
			wantBody:             `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 400,
		},
		{
			name:                 "err, id is negative",
			id:                   "-1",
			body:                 `{"title":"New title","author":"New Author","year":2000,"available":false}`,
			serviceErr:           nil,
			serviceBook:          model.Book{},
			wantServiceCallCount: 0,
			wantBody:             `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 400,
		},
		{
			name:                 "err, id is not number",
			id:                   "1lgdf",
			body:                 `{"title":"New title","author":"New Author","year":2000,"available":false}`,
			serviceErr:           nil,
			serviceBook:          model.Book{},
			wantServiceCallCount: 0,
			wantBody:             `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 400,
		},
		{
			name:                 "err, bad json",
			id:                   "1",
			body:                 `{"title":"New title","author":"New Author","year":2000,"available":}`,
			serviceErr:           nil,
			serviceBook:          model.Book{},
			wantServiceCallCount: 0,
			wantBody:             `{"error":"invalid character '}' looking for beginning of value"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 400,
		},
		{
			name:                 "err, year is string",
			id:                   "1",
			body:                 `{"title":"New title","author":"New Author","year":"2000","available":false}`,
			serviceErr:           nil,
			serviceBook:          model.Book{},
			wantServiceCallCount: 0,
			wantBody:             `{"error":"json: cannot unmarshal string into Go struct field EditBookBody.year of type uint16"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 400,
		},
		{
			name:                 "err, available is string",
			id:                   "1",
			body:                 `{"title":"New title","author":"New Author","year":2000,"available":"false"}`,
			serviceErr:           nil,
			serviceBook:          model.Book{},
			wantServiceCallCount: 0,
			wantBody:             `{"error":"json: cannot unmarshal string into Go struct field EditBookBody.available of type bool"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 400,
		},
		{
			name:                 "err, unknown fiels",
			id:                   "1",
			body:                 `{"title":"New title","author":"New Author","year":2000,"available":false,"test":10}`,
			serviceErr:           nil,
			serviceBook:          model.Book{},
			wantServiceCallCount: 0,
			wantBody:             `{"error":"json: unknown field \"test\""}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 400,
		},
		{
			name:                 "err, not found",
			id:                   "1",
			body:                 `{"title":"New title","author":"New Author","year":2000,"available":false}`,
			serviceErr:           apperrors.ErrNotFound,
			serviceBook:          model.Book{},
			wantServiceCallCount: 1,
			wantId:               1,
			wantBody:             `{"error":"` + apperrors.ErrNotFound.Error() + `"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 404,
		},
		{
			name:                 "err, empty json",
			id:                   "1",
			body:                 `{}`,
			serviceErr:           apperrors.ErrRequiredFields,
			serviceBook:          model.Book{},
			wantServiceCallCount: 1,
			wantId:               1,
			wantBody:             `{"error":"` + apperrors.ErrRequiredFields.Error() + `"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 400,
		},
		{
			name:                 "err, empty json",
			id:                   "1",
			body:                 `{"title":"New title","author":"New Author","year":2100,"available":false}`,
			serviceErr:           apperrors.ErrInvalidValues,
			serviceBook:          model.Book{},
			wantServiceCallCount: 1,
			wantId:               1,
			wantBody:             `{"error":"` + apperrors.ErrInvalidValues.Error() + `"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 400,
		},
		{
			name:                 "err, cancel in service",
			id:                   "1",
			body:                 `{"title":"New title","author":"New Author","year":2000,"available":false}`,
			serviceErr:           context.Canceled,
			serviceBook:          model.Book{},
			wantServiceCallCount: 1,
			wantId:               1,
			wantBody:             ``,
			wantHeaders:          map[string][]string{},
			wantCode:             200,
		},
		{
			name:                 "err, deadline exceeded in service",
			id:                   "1",
			body:                 `{"title":"New title","author":"New Author","year":2000,"available":false}`,
			serviceErr:           context.DeadlineExceeded,
			serviceBook:          model.Book{},
			wantServiceCallCount: 1,
			wantId:               1,
			wantBody:             `{"error":"context deadline exceeded"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: 504,
		},
		{
			name:                 "err, cancel in service",
			id:                   "1",
			cancel:               true,
			body:                 `{"title":"New title","author":"New Author","year":2000,"available":false}`,
			wantServiceCallCount: 0,
			wantBody:             ``,
			wantHeaders:          map[string][]string{},
			wantCode:             200,
		},
		{
			name:                 "err, cancel in service",
			id:                   "1",
			deadlineContext:      true,
			body:                 `{"title":"New title","author":"New Author","year":2000,"available":false}`,
			wantServiceCallCount: 0,
			wantBody:             ``,
			wantHeaders:          map[string][]string{},
			wantCode:             504,
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
				http.MethodPut,
				"http://localhost:8080/books/"+tt.id,
				strings.NewReader(tt.body),
			)

			rr := httptest.NewRecorder()

			s := ServiceMock{
				book: tt.serviceBook,
				err:  tt.serviceErr,
			}

			if tt.cancel {
				cancel()
			}

			req.SetPathValue("id", tt.id)
			handler.EditBookByID(&s)(rr, req)

			if ctx != s.ctx && tt.wantServiceCallCount >= 1 {
				t.Error("context passed to service is not equal to request's context")
			}

			if !compareHeaders(tt.wantHeaders, rr.Header()) {
				t.Errorf(
					"expected headers to be %v\nbut got %v",
					tt.wantHeaders,
					rr.Header(),
				)
			}

			if tt.wantCode != rr.Code {
				t.Errorf(
					"expected code to be %d\nbut got %d",
					tt.wantCode,
					rr.Code,
				)
			}

			if tt.wantServiceCallCount != s.callCount {
				t.Errorf(
					"expected service to be called %d times but got %d",
					tt.wantServiceCallCount,
					s.callCount,
				)
			}

			if tt.wantId != s.id && tt.wantServiceCallCount >= 1 {
				t.Errorf(
					"expected service to be called with id %d but got %d",
					tt.wantId,
					s.id,
				)
			}

			if tt.wantServiceCallCount >= 1 {
				var expectedInput service.EditBookInput
				decoder := json.NewDecoder(strings.NewReader(tt.body))
				err := decoder.Decode(&expectedInput)
				if err != nil {
					t.Fatal("failed to decode body")
				}
				if !reflect.DeepEqual(expectedInput, s.editInput) {
					t.Errorf(
						"expected to call service with input %v\nbut got %v",
						expectedInput,
						s.editInput,
					)
				}
			}

			if strings.TrimSpace(tt.wantBody) != strings.TrimSpace(rr.Body.String()) {
				t.Errorf(
					"expected response body to be %v\nbut got %v",
					strings.TrimSpace(tt.wantBody),
					strings.TrimSpace(rr.Body.String()),
				)
			}
		})
	}
}
