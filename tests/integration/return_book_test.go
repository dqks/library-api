package integration

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/handler"
	"library-api/internal/service"
	"library-api/internal/testutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestReturnBook(t *testing.T) {
	tests := []struct {
		// init
		name            string
		body            string
		loanID          string
		deadlineContext bool
		cancel          bool
		// assert return book
		wantReturnBody string
		wantReturnCode int
		wantHeaders    http.Header
		// assert get book
		wantGetBookBody string
		// assert get loan
		wantGetLoanBody string
	}{
		{
			name:           "success, with body",
			body:           `{"return_date":"2026-10-29"}`,
			loanID:         "1",
			wantReturnBody: `{"success":true}`,
			wantReturnCode: 200,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBookBody: `{"id":1,"title":"Title 1","author":"Author 1","year":2000,"available":true}`,
			wantGetLoanBody: `{"id":1,"book_id":1,"user_id":1,"loan_date":"2026-10-29","return_date":"2026-10-29"}`,
		},
		{
			name:           "success, no body",
			body:           ``,
			loanID:         "1",
			wantReturnBody: `{"success":true}`,
			wantReturnCode: 200,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBookBody: `{"id":1,"title":"Title 1","author":"Author 1","year":2000,"available":true}`,
			wantGetLoanBody: `{"id":1,"book_id":1,"user_id":1,"loan_date":"2026-10-29","return_date":"` + time.Now().Format(time.DateOnly) + `"}`,
		},
		{
			name:           "canceled",
			cancel:         true,
			body:           ``,
			loanID:         "1",
			wantReturnCode: 200,
		},
		{
			name:            "deadline exceeded",
			deadlineContext: true,
			body:            ``,
			loanID:          "1",
			wantReturnCode:  apperrors.ErrDeadlineExceededCode,
		},
		{
			name:           "err, loan not found",
			body:           `{"return_date":"2026-10-29"}`,
			loanID:         "55",
			wantReturnBody: `{"error":"` + apperrors.ErrNotFound.Error() + `"}`,
			wantReturnCode: 404,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBookBody: `{"id":1,"title":"Title 1","author":"Author 1","year":2000,"available":true}`,
			wantGetLoanBody: `{"error":"` + apperrors.ErrNotFound.Error() + `"}`,
		},
		{
			name:           "err, bad path value",
			body:           `{"return_date":"2026-10-29"}`,
			loanID:         "1fasdf",
			wantReturnBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			wantReturnCode: 400,
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantGetBookBody: `{"id":1,"title":"Title 1","author":"Author 1","year":2000,"available":true}`,
			wantGetLoanBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, pool := testutil.CreatePostgresTestRepo(t, "../../.env")
			t.Cleanup(func() {
				pool.Close()
			})
			serviceLoan := service.CreateLoanService(repo)
			serviceBooks := service.CreateBookService(repo)

			testutil.SeedBooks(t, pool)
			testutil.SeedUsers(t, pool)
			testutil.SeedLoans(t, pool)

			t.Cleanup(func() {
				ctx := context.Background()
				_, err := testutil.ClearLoans(ctx, pool)
				if err != nil {
					t.Errorf("failed to clear loans %v", err)
				}
				_, err = testutil.ClearBooks(ctx, pool)
				if err != nil {
					t.Errorf("failed to clear books %v", err)
				}
				_, err = testutil.ClearUsers(ctx, pool)
				if err != nil {
					t.Errorf("failed to clear users %v", err)
				}
			})

			mux := http.NewServeMux()
			mux.HandleFunc("POST /loans/{loanId}", handler.ReturnBook(serviceLoan))
			mux.HandleFunc("GET /books/{id}", handler.GetBookByID(serviceBooks))
			mux.HandleFunc("GET /loans/{id}", handler.GetLoan(serviceLoan))

			//Return book request
			var ctxReturn context.Context
			var cancelReturn context.CancelFunc

			if tt.deadlineContext {
				ctxReturn, cancelReturn = context.WithDeadline(context.Background(), time.Time{})
				defer cancelReturn()
			} else {
				ctxReturn, cancelReturn = context.WithCancel(context.Background())
				defer cancelReturn()
			}

			if tt.cancel {
				cancelReturn()
			}

			reqReturn := httptest.NewRequestWithContext(
				ctxReturn,
				http.MethodPost,
				"/loans/"+tt.loanID,
				strings.NewReader(tt.body),
			)
			rr := httptest.NewRecorder()

			if tt.body != "" {
				reqReturn.Header.Set("Content-Type", "application/json")

			}

			mux.ServeHTTP(rr, reqReturn)

			if rr.Code != tt.wantReturnCode {
				t.Errorf(
					"expected code to be %d but got %d",
					tt.wantReturnCode,
					rr.Code,
				)
			}

			if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(tt.wantReturnBody) {
				t.Errorf(
					"expected return book body to be %v\nbut got %v",
					tt.wantReturnBody,
					rr.Body.String(),
				)
			}

			if !testutil.CompareHeaders(rr.Header(), tt.wantHeaders) {
				t.Errorf(
					"expected headers to be %v\nbut got %v",
					tt.wantHeaders,
					rr.Header(),
				)
			}

			// Get book

			if tt.cancel || tt.deadlineContext {
				return
			}

			ctxGet := context.Background()
			rrGet := httptest.NewRecorder()
			reqGet := httptest.NewRequestWithContext(
				ctxGet,
				http.MethodGet,
				`/books/1`,
				nil,
			)

			mux.ServeHTTP(rrGet, reqGet)

			if strings.TrimSpace(rrGet.Body.String()) != strings.TrimSpace(tt.wantGetBookBody) {
				t.Errorf(
					"expected get book body to be %v\nbut got %v",
					tt.wantGetBookBody,
					rrGet.Body.String(),
				)
			}

			// Get loan

			ctxLoan := context.Background()
			rrLoan := httptest.NewRecorder()
			reqLoan := httptest.NewRequestWithContext(
				ctxLoan,
				http.MethodGet,
				`/loans/`+tt.loanID,
				nil,
			)

			mux.ServeHTTP(rrLoan, reqLoan)

			if strings.TrimSpace(rrLoan.Body.String()) != strings.TrimSpace(tt.wantGetLoanBody) {
				t.Errorf(
					"expected get loan body to be %v\nbut got %v",
					tt.wantGetLoanBody,
					rrLoan.Body.String(),
				)
			}
		})
	}
}
