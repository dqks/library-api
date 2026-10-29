package tests

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/handler"
	"library-api/internal/model"
	"library-api/internal/testutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestGetUser(t *testing.T) {
	tests := []struct {
		name            string
		id              string
		cancel          bool
		deadlineContext bool
		// service mock
		serviceUser model.User
		serviceErr  error
		// asserts
		wantBody        string
		wantCode        int
		wantHeaders     http.Header
		wantServiceCall int
	}{
		{
			name: "success",
			id:   "1",
			//service mock
			serviceUser: model.User{
				ID:   1,
				Name: "test",
			},
			serviceErr: nil,
			//asserts
			wantBody: `{"id":1,"name":"test"}`,
			wantCode: 200,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantServiceCall: 1,
		},
		{
			name: "zero path value",
			id:   "0",
			//service mock
			serviceErr: nil,
			//asserts
			wantBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantCode: 400,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
		},
		{
			name: "negative path value",
			id:   "-1",
			//service mock
			serviceErr: nil,
			//asserts
			wantBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantCode: 400,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
		},
		{
			name: "service canceled",
			id:   "1",
			//service mock
			serviceErr: context.Canceled,
			//asserts
			wantCode:        200,
			wantServiceCall: 1,
		},
		{
			name: "service deadline",
			id:   "1",
			//service mock
			serviceErr: context.DeadlineExceeded,
			//asserts
			wantCode:        apperrors.ErrDeadlineExceededCode,
			wantBody:        `{"error":"` + context.DeadlineExceeded.Error() + `"}`,
			wantServiceCall: 1,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
		},
		{
			name: "service user not found",
			id:   "1",
			//service mock
			serviceErr: apperrors.ErrNotFound,
			//asserts
			wantCode:        apperrors.ErrNotFoundCode,
			wantBody:        `{"error":"` + apperrors.ErrNotFound.Error() + `"}`,
			wantServiceCall: 1,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
		},
		{
			name: "service internal",
			id:   "1",
			//service mock
			serviceErr: apperrors.ErrInternal,
			//asserts
			wantCode:        apperrors.ErrInternalCode,
			wantBody:        `{"error":"` + apperrors.ErrInternal.Error() + `"}`,
			wantServiceCall: 1,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
		},
		{
			name:     "handler canceled",
			id:       "1",
			cancel:   true,
			wantCode: 200,
		},
		{
			name:            "handler deadline",
			id:              "1",
			deadlineContext: true,
			wantCode:        apperrors.ErrDeadlineExceededCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := testutil.GetNeededCtx(tt.deadlineContext)
			t.Cleanup(func() {
				cancel()
			})
			if tt.cancel {
				cancel()
			}

			rr := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(
				ctx,
				http.MethodGet,
				`http://localhost:8080/users/`+tt.id,
				nil,
			)
			req.SetPathValue("id", tt.id)

			s := ServiceMock{
				user: tt.serviceUser,
				err:  tt.serviceErr,
			}

			handler.GetUser(&s)(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf(
					"expected code %d but got %d",
					tt.wantCode,
					rr.Code,
				)
			}

			if s.callCount != tt.wantServiceCall {
				t.Errorf(
					"expected to call service %d times but got %d",
					tt.wantServiceCall,
					s.callCount,
				)
			}

			if s.callCount >= 1 {
				id, err := strconv.Atoi(tt.id)

				if err != nil {
					t.Errorf("failed to convert id to int")
				}

				if s.id != id {
					t.Errorf(
						"expected to call service %d times but got %d",
						tt.wantServiceCall,
						s.callCount,
					)
				}

				if s.ctx != ctx {
					t.Errorf(
						"expected to call service wtih ctx %v  but got %v",
						ctx,
						s.ctx,
					)
				}
			}

			if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(tt.wantBody) {
				t.Errorf(
					"expected body %v\nbut got %v",
					tt.wantBody,
					rr.Body.String(),
				)
			}

			if !testutil.CompareHeaders(rr.Header(), tt.wantHeaders) {
				t.Errorf(
					"expected headers %v\nbut got %v",
					tt.wantHeaders,
					rr.Header(),
				)
			}

		})
	}
}
