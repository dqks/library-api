package tests

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/handler"
	"library-api/internal/model"
	"library-api/internal/testutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name            string
		cancel          bool
		deadlineContext bool
		body            string
		//service mock
		serviceUser model.User
		serviceErr  error
		//assetions
		wantServiceCallCount int
		wantBody             string
		wantCode             int
		wantHeaders          http.Header
	}{
		{
			name: "success",
			body: `{"name":"test"}`,
			//service mock
			serviceUser: model.User{
				ID:   1,
				Name: "test",
			},
			serviceErr: nil,
			//assetions
			wantServiceCallCount: 1,
			wantBody:             `{"id":1,"name":"test"}`,
			wantCode:             201,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
		},
		{
			name: "bad body",
			body: `{"name":"te}`,
			//assetions
			wantServiceCallCount: 0,
			wantBody:             `{"error":"unexpected EOF"}`,
			wantCode:             400,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
		},
		{
			name: "bad field type",
			body: `{"name":143}`,
			//assetions
			wantServiceCallCount: 0,
			wantBody:             `{"error":"json: cannot unmarshal number into Go struct field .name of type string"}`,
			wantCode:             400,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
		},
		{
			name: "no needed field body",
			body: `{"namewrong":143}`,
			//assetions
			wantServiceCallCount: 0,
			wantBody:             `{"error":"json: unknown field \"namewrong\""}`,
			wantCode:             400,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
		},
		{
			name: "additional field in body",
			body: `{"name":"test","add":"das"}`,
			//assetions
			wantServiceCallCount: 0,
			wantBody:             `{"error":"json: unknown field \"add\""}`,
			wantCode:             400,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
		},
		{
			name: "service canceled",
			body: `{"name":"test"}`,
			// service mock
			serviceErr: context.Canceled,
			//assetions
			wantServiceCallCount: 1,
			wantCode:             200,
		},
		{
			name: "service deadline exceeded",
			body: `{"name":"test"}`,
			// service mock
			serviceErr: context.DeadlineExceeded,
			//assetions
			wantServiceCallCount: 1,
			wantBody:             `{"error":"context deadline exceeded"}`,
			wantCode:             apperrors.ErrDeadlineExceededCode,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
		},
		{
			name: "service internal err",
			body: `{"name":"test"}`,
			// service mock
			serviceErr: apperrors.ErrInternal,
			//assetions
			wantServiceCallCount: 1,
			wantBody:             `{"error":"` + apperrors.ErrInternal.Error() + `"}`,
			wantCode:             apperrors.ErrInternalCode,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := testutil.GetNeededCtx(tt.deadlineContext)
			t.Cleanup(func() {
				cancel()
			})

			rr := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(
				ctx,
				http.MethodPost,
				"http://localhost:8080/users",
				strings.NewReader(tt.body),
			)

			s := ServiceMock{
				user: tt.serviceUser,
				err:  tt.serviceErr,
			}

			handler.CreateUser(&s)(rr, req)

			if s.callCount != tt.wantServiceCallCount {
				t.Errorf(
					"expected to call serivce %d times but got %d",
					tt.wantServiceCallCount,
					s.callCount,
				)
			}

			if len(strings.TrimSpace(rr.Body.String())) != len(strings.TrimSpace(tt.wantBody)) {
				t.Errorf(
					"expected body %v\nbut got %v",
					tt.wantBody,
					rr.Body.String(),
				)
			}

			if rr.Code != tt.wantCode {
				t.Errorf(
					"expected code %d but got %d",
					tt.wantCode,
					rr.Code,
				)
			}

			if !testutil.CompareHeaders(rr.Header(), tt.wantHeaders) {
				t.Errorf(
					"expected headers %v\nbut got %v",
					tt.wantHeaders,
					rr.Header(),
				)
			}

			if s.ctx != ctx && s.callCount >= 1 {
				t.Errorf(
					"expected to call service with context %v\nbut got %v",
					ctx,
					s.ctx,
				)
			}
		})
	}
}
