package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"library-api/internal/apperrors"
	"net/http"
)

type loanBookResponse struct {
	Success bool `json:"success"`
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
			bookID, err := ValidatePathID(w, r.PathValue("bookId"))
			if err != nil {
				return
			}

			userID, err := ValidatePathID(w, r.PathValue("userId"))
			if err != nil {
				return
			}

			err = s.LoanBook(ctx, bookID, userID)

			if err != nil {
				errors := []error{
					context.Canceled,
					context.DeadlineExceeded,
					apperrors.ErrNotFound,
				}
				ProcessServiceErr(w, err, errors)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(201)
			encoder := json.NewEncoder(w)
			if err := encoder.Encode(loanBookResponse{Success: true}); err != nil {
				fmt.Println(apperrors.ErrInternal.Error())
			}
		}
	}
}
