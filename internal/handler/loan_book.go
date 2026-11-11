package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"library-api/internal/apperrors"
	"log/slog"
	"net/http"
	"time"
)

type LoanBookResponse struct {
	ID       int    `json:"id"`
	LoanDate string `json:"loan_date"`
}

func LoanBook(s LoanService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				w.WriteHeader(apperrors.ErrDeadlineExceededCode)
			}
		default:
			bookID, err := GetPathID(w, r.PathValue("bookId"))
			if err != nil {
				return
			}

			userID, err := GetPathID(w, r.PathValue("userId"))
			if err != nil {
				return
			}

			loan, err := s.LoanBook(ctx, bookID, userID)

			if err != nil {
				errors := []error{
					context.Canceled,
					context.DeadlineExceeded,
					apperrors.ErrNotFound,
					apperrors.ErrBookUnavailable,
					apperrors.ErrUserLoanLimit,
				}
				ProcessServiceErr(w, err, errors)
				slog.Error(
					"loan book error",
					"user_id", userID,
					"bookID", bookID,
					"err", err,
				)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(201)
			encoder := json.NewEncoder(w)
			response := LoanBookResponse{
				ID:       loan.ID,
				LoanDate: loan.LoanDate.Format(time.DateOnly),
			}

			if err := encoder.Encode(response); err != nil {
				fmt.Println(apperrors.ErrInternal.Error())
			}
		}
	}
}
