package integration

import (
	"context"
	"encoding/json"
	"library-api/internal/apperrors"
	"library-api/internal/database"
	"library-api/internal/handler"
	"library-api/internal/service"
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
		name            string
		bookID          string
		userID          string
		cancel          bool
		deadlineContext bool
		body            string
		secondLoan      bool
		// asserts loan
		wantLoanCode   int
		wantLoanHeader http.Header
		wantLoanBody   string
		// asserts get book
		wantGetBookBody string
		// asserts get loan
		wantGetLoanBody string
	}{
		{
			name:   "success",
			bookID: "1",
			userID: "1",
			// asserts loan
			wantLoanCode: 201,
			wantLoanHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantLoanBody: `{"id":3,"loan_date":"2026-11-02"}`,
			// asserts get book
			wantGetBookBody: `{"id":1,"title":"Title 1","author":"Author 1","year":2000,"available":false}`,
			//assert get loan
			wantGetLoanBody: `{"id":3,"book_id":1,"user_id":1,"loan_date":"2026-11-02","return_date":null}`,
		},
		{
			name:   "err, book not found",
			bookID: "100",
			userID: "1",
			// asserts loan
			wantLoanCode: apperrors.ErrNotFoundCode,
			wantLoanHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantLoanBody: `{"error":"` + apperrors.ErrNotFound.Error() + `"}`,
			// asserts get book
			wantGetBookBody: `{"error":"` + apperrors.ErrNotFound.Error() + `"}`,
			//assert get loan
		},
		{
			name:   "user, not found",
			bookID: "1",
			userID: "100",
			// asserts loan
			wantLoanCode: apperrors.ErrNotFoundCode,
			wantLoanHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantLoanBody: `{"error":"` + apperrors.ErrNotFound.Error() + `"}`,
			// asserts get book
			wantGetBookBody: `{"id":1,"title":"Title 1","author":"Author 1","year":2000,"available":true}`,
			//assert get loan
		},
		{
			name:   "canceled",
			bookID: "1",
			userID: "1",
			cancel: true,
			// asserts loan
			wantLoanCode: 200,
		},
		{
			name:            "deadline exceeded",
			bookID:          "1",
			userID:          "1",
			deadlineContext: true,
			// asserts loan
			wantLoanCode: apperrors.ErrDeadlineExceededCode,
		},
		{
			name:   "err, bad path value for book id",
			bookID: "1gdnkghdn",
			userID: "1",
			// asserts loan
			wantLoanCode: apperrors.ErrIncorrectPathValueCode,
			wantLoanHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantLoanBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			// asserts get book
			wantGetBookBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			//assert get loan
		},
		{
			name:   "err, bad path value for user id",
			bookID: "1",
			userID: "1gsdgkdf",
			// asserts loan
			wantLoanCode: apperrors.ErrIncorrectPathValueCode,
			wantLoanHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantLoanBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			// asserts get book
			wantGetBookBody: `{"id":1,"title":"Title 1","author":"Author 1","year":2000,"available":true}`,
			//assert get loan
		},
		{
			name:   "err, zero path value for book id",
			bookID: "0",
			userID: "1",
			// asserts loan
			wantLoanCode: apperrors.ErrIncorrectPathValueCode,
			wantLoanHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantLoanBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			// asserts get book
			wantGetBookBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			//assert get loan
		},
		{
			name:   "err, zero path value for user id",
			bookID: "1",
			userID: "0",
			// asserts loan
			wantLoanCode: apperrors.ErrIncorrectPathValueCode,
			wantLoanHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantLoanBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			// asserts get book
			wantGetBookBody: `{"id":1,"title":"Title 1","author":"Author 1","year":2000,"available":true}`,
			//assert get loan
		},
		{
			name:   "err, negative path value for book id",
			bookID: "-1",
			userID: "1",
			// asserts loan
			wantLoanCode: apperrors.ErrIncorrectPathValueCode,
			wantLoanHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantLoanBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			// asserts get book
			wantGetBookBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			//assert get loan
		},
		{
			name:   "err, negative path value for user id",
			bookID: "1",
			userID: "-1",
			// asserts loan
			wantLoanCode: apperrors.ErrIncorrectPathValueCode,
			wantLoanHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantLoanBody: `{"error":"` + apperrors.ErrIncorrectPathValue.Error() + `"}`,
			// asserts get book
			wantGetBookBody: `{"id":1,"title":"Title 1","author":"Author 1","year":2000,"available":true}`,
			//assert get loan
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn, err := database.GetDsn("../../.env", "DATABASE_TEST_URL")

			if err != nil {
				t.Fatalf("ошибка загрузки dsn: %v", err)
			}

			repo, pool := testutil.CreatePostgresTestRepo(t, dsn)
			t.Cleanup(func() {
				pool.Close()
			})

			err = database.UpMigrations(pool, "../../migrations")

			if err != nil {
				t.Fatalf("ошибка применения миграций: %v", err)
			}

			t.Cleanup(func() {
				err := database.DownMigrations(pool, "../../migrations")
				if err != nil {
					t.Fatalf("ошибка применения миграций: %v", err)
				}
			})

			serviceBook := service.CreateBookService(repo)
			serviceLoan := service.CreateLoanService(repo)

			testutil.SeedBooks(t, pool)
			testutil.SeedUsers(t, pool)
			testutil.SeedLoans(t, pool)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /loans/{bookId}/{userId}", handler.LoanBook(serviceLoan))
			mux.HandleFunc("GET /books/{id}", handler.GetBookByID(serviceBook))
			mux.HandleFunc("GET /loans/{id}", handler.GetLoan(serviceLoan))

			// Loan Book
			var ctxLoan context.Context
			var cancel context.CancelFunc

			if tt.deadlineContext {
				ctxLoan, cancel = context.WithDeadline(context.Background(), time.Time{})
				defer cancel()
			} else {
				ctxLoan, cancel = context.WithCancel(context.Background())
				defer cancel()
			}

			if tt.cancel {
				cancel()
			}

			rrLoan := httptest.NewRecorder()
			reqLoan := httptest.NewRequestWithContext(
				ctxLoan,
				http.MethodPost,
				"/loans/"+tt.bookID+"/"+tt.userID,
				strings.NewReader(tt.body),
			)

			mux.ServeHTTP(rrLoan, reqLoan)

			if strings.TrimSpace(rrLoan.Body.String()) != strings.TrimSpace(tt.wantLoanBody) {
				t.Errorf(
					"expected body to be %v be\nbut got %v",
					tt.wantLoanBody,
					rrLoan.Body.String(),
				)
			}

			if !testutil.CompareHeaders(tt.wantLoanHeader, rrLoan.Header()) {
				t.Errorf(
					"expected header to be %v\nbut got %v",
					tt.wantLoanHeader,
					rrLoan.Header(),
				)
			}

			if tt.wantLoanCode != rrLoan.Code {
				t.Errorf(
					"expected code to be %d but got %d",
					tt.wantLoanCode,
					rrLoan.Code,
				)
			}

			if !tt.cancel && !tt.deadlineContext {
				// Get Book
				ctxGetBook := context.Background()
				rrGetBook := httptest.NewRecorder()
				reqGetBook := httptest.NewRequestWithContext(
					ctxGetBook,
					http.MethodGet,
					`/books/`+tt.bookID,
					nil,
				)

				mux.ServeHTTP(rrGetBook, reqGetBook)

				if strings.TrimSpace(rrGetBook.Body.String()) != strings.TrimSpace(tt.wantGetBookBody) {
					t.Errorf(
						"expected get book body to be %v\nbut got %v",
						tt.wantGetBookBody,
						rrGetBook.Body.String(),
					)
				}

				// Get Loan
				var loanResponse handler.LoanBookResponse
				decoder := json.NewDecoder(rrLoan.Body)
				if err = decoder.Decode(&loanResponse); err != nil {
					t.Fatalf("failed to decode loan book body %v", err)
				}

				if loanResponse.ID != 0 {
					ctxGetLoan := context.Background()
					rrGetLoan := httptest.NewRecorder()
					reqGetLoan := httptest.NewRequestWithContext(
						ctxGetLoan,
						http.MethodGet,
						`/loans/`+strconv.Itoa(loanResponse.ID),
						nil,
					)

					mux.ServeHTTP(rrGetLoan, reqGetLoan)

					if strings.TrimSpace(rrGetLoan.Body.String()) != strings.TrimSpace(tt.wantGetLoanBody) {
						t.Errorf(
							"expected get loan body to be %v be\nbut got %v",
							tt.wantGetLoanBody,
							rrGetLoan.Body.String(),
						)
					}
				}
			}
		})
	}
}
