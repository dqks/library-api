package tests

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/handler"
	"library-api/internal/testutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestLoanBook(t *testing.T) {
	tests := []struct {
		name                 string
		cancel               bool
		deadlineExceeded     bool
		bookID               string
		userID               string
		loanErr              error
		wantHeader           http.Header
		wantBody             string
		wantCode             int
		wantServiceCallCount int
	}{
		{
			name:    "success",
			userID:  "1",
			bookID:  "1",
			loanErr: nil,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantBody:             `{"success":true}`,
			wantCode:             201,
			wantServiceCallCount: 1,
		},
		{
			name:     "canceled",
			cancel:   true,
			userID:   "1",
			bookID:   "1",
			loanErr:  nil,
			wantCode: 200,
		},
		{
			name:             "deadline exceeded",
			deadlineExceeded: true,
			userID:           "1",
			bookID:           "1",
			loanErr:          nil,
			wantCode:         504,
		},
		{
			name:    "err, not found",
			userID:  "1",
			bookID:  "1",
			loanErr: apperrors.ErrNotFound,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantBody:             `{"error":"` + apperrors.ErrNotFound.Error() + `"}`,
			wantCode:             404,
			wantServiceCallCount: 1,
		},
		{
			name:    "err, bad bookID",
			userID:  "1",
			bookID:  "1gdf",
			loanErr: nil,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantCode: 400,
		},
		{
			name:    "err, zero bookID",
			userID:  "1",
			bookID:  "0",
			loanErr: nil,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantCode: 400,
		},
		{
			name:    "err, negative bookID",
			userID:  "1",
			bookID:  "-1",
			loanErr: nil,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantCode: 400,
		},
		{
			name:    "err, bad userID",
			userID:  "1das",
			bookID:  "1",
			loanErr: nil,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantCode: 400,
		},
		{
			name:    "err, zero userID",
			userID:  "0",
			bookID:  "1",
			loanErr: nil,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantCode: 400,
		},
		{
			name:    "err, negative bookID",
			userID:  "-1",
			bookID:  "1",
			loanErr: nil,
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx context.Context
			var cancel context.CancelFunc
			if tt.deadlineExceeded {
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
				http.MethodPost,
				"http://localhost:8080/loans/"+tt.bookID+"/"+tt.userID,
				nil,
			)
			req.SetPathValue("bookId", tt.bookID)
			req.SetPathValue("userId", tt.userID)
			req.Header.Set("Content-Type", "application/json")

			s := ServiceMock{
				err: tt.loanErr,
			}

			handler.LoanBook(&s)(rr, req)

			if ctx != s.ctx && tt.wantServiceCallCount >= 1 {
				t.Errorf(
					"expected to call handler with ctx be %v\nbut got %v",
					ctx,
					s.ctx,
				)
			}

			if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(tt.wantBody) {
				t.Errorf(
					"expected body to be %v\nbut got %v",
					tt.wantBody,
					rr.Body.String(),
				)
			}

			if tt.wantServiceCallCount != s.callCount {
				t.Errorf(
					"expected service call count to be %d\nbut got %d",
					tt.wantServiceCallCount,
					s.callCount,
				)
			}

			if tt.wantServiceCallCount >= 1 {
				bookID, err := strconv.Atoi(tt.bookID)

				if err != nil {
					t.Fatalf("failed to convert string to int")
				}

				userID, err := strconv.Atoi(tt.userID)

				if err != nil {
					t.Fatalf("failed to convert string to int")
				}

				if bookID != s.bookID {
					t.Errorf(
						"expected to call service with bookID  %v\nbut got %d",
						tt.bookID,
						s.bookID,
					)
				}

				if userID != s.userID {
					t.Errorf(
						"expected to call service with userID  %v\nbut got %d",
						tt.userID,
						s.userID,
					)
				}
			}

			if rr.Code != tt.wantCode {
				t.Errorf(
					"expected code to be %d\nbut got %d",
					tt.wantCode,
					rr.Code,
				)
			}

			if !testutil.CompareHeaders(rr.Header(), tt.wantHeader) {
				t.Errorf(
					"expected headers to be %v\nbut got %v",
					tt.wantHeader,
					rr.Header(),
				)
			}
		})
	}
}
