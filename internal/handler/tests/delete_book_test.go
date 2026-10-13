package tests

import (
	"context"
	"errors"
	"library-api/internal/apperrors"
	"library-api/internal/handler"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestDeleteBookByID(t *testing.T) {
	tests := []struct {
		name                 string
		id                   int
		deadlineContext      bool
		cancel               bool
		wantServiceCallCount int
		serviceErr           error
		wantCode             int
		wantBody             string
		wantHeaders          http.Header
	}{
		{
			name:       "success",
			id:         1,
			serviceErr: nil,
			wantCode:   200,
			wantBody:   `{"success":true}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantServiceCallCount: 1,
		},
		{
			name:       "err, id is less than zero",
			id:         -1,
			serviceErr: nil,
			wantCode:   400,
			wantBody:   `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantServiceCallCount: 0,
		},
		{
			name:       "err, id is zero",
			id:         0,
			serviceErr: nil,
			wantCode:   400,
			wantBody:   `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantServiceCallCount: 0,
		},
		{
			name:       "err, not found",
			id:         1,
			serviceErr: apperrors.ErrNotFound,
			wantCode:   404,
			wantBody:   `{"error":"` + apperrors.ErrNotFound.Error() + `"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantServiceCallCount: 1,
		},
		{
			name:                 "err, context is canceled in service",
			id:                   1,
			serviceErr:           context.Canceled,
			wantCode:             200,
			wantBody:             ``,
			wantHeaders:          map[string][]string{},
			wantServiceCallCount: 1,
		},
		{
			name:       "err, context deadline exceeded in service",
			id:         1,
			serviceErr: context.DeadlineExceeded,
			wantCode:   504,
			wantBody:   `{"error":"` + context.DeadlineExceeded.Error() + `"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantServiceCallCount: 1,
		},
		{
			name:                 "cancel in repository",
			id:                   1,
			serviceErr:           nil,
			wantCode:             200,
			cancel:               true,
			wantBody:             ``,
			wantHeaders:          map[string][]string{},
			wantServiceCallCount: 0,
		},
		{
			name:                 "deadline exceeded in repository",
			id:                   1,
			serviceErr:           nil,
			wantCode:             504,
			deadlineContext:      true,
			wantBody:             ``,
			wantHeaders:          map[string][]string{},
			wantServiceCallCount: 0,
		},
		{
			name:       "unknown error in service",
			id:         1,
			serviceErr: errors.New("unknown error"),
			wantCode:   500,
			wantBody:   `{"error":"unknown error"}`,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantServiceCallCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx context.Context
			var cancel context.CancelFunc

			if !tt.deadlineContext {
				ctx, cancel = context.WithCancel(context.Background())
				defer cancel()
			} else {
				ctx, cancel = context.WithDeadline(context.Background(), time.Time{})
				defer cancel()
			}

			rr := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(
				ctx,
				http.MethodDelete,
				"http://localhost:8080/books/"+strconv.Itoa(tt.id),
				nil,
			)
			req.SetPathValue("id", strconv.Itoa(tt.id))

			s := ServiceMock{
				err: tt.serviceErr,
			}

			if tt.cancel {
				cancel()
			}

			handler.DeleteBookByID(&s)(rr, req)

			if ctx != s.ctx && tt.wantServiceCallCount >= 1 {
				t.Fatalf("expected context to be %v but got %v", ctx, s.ctx)
			}

			if s.id != tt.id && tt.wantServiceCallCount >= 1 {
				t.Fatalf("expeceted to pass to service id %d but got %d", tt.id, s.id)
			}

			if rr.Code != tt.wantCode {
				t.Fatalf("expeceted code to be %d but got %d", tt.wantCode, rr.Code)
			}

			if s.callCount != tt.wantServiceCallCount {
				t.Fatalf(
					"expected service to be called %d times but got %d",
					tt.wantServiceCallCount,
					s.callCount,
				)
			}

			if !compareHeaders(tt.wantHeaders, rr.Header()) {
				t.Fatalf(
					"expected headers to be %v\nbut got %v",
					tt.wantHeaders,
					rr.Header(),
				)
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
