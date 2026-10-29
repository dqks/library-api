package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"library-api/internal/apperrors"
	"net/http"
)

type ReturnBookResponse struct {
	Success bool `json:"success"`
}

type ReturnBookBody struct {
	ReturnDate *string `json:"return_date"`
}

func ReturnBook(s LoanService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				w.WriteHeader(apperrors.ErrDeadlineExceededCode)
			}
		default:
			loanID, err := ValidatePathID(w, r.PathValue("loanId"))

			if err != nil {
				return
			}

			var body ReturnBookBody
			decoder := json.NewDecoder(r.Body)
			if err = decoder.Decode(&body); err != nil {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(400)
				encoder := json.NewEncoder(w)
				if encodeErr := encoder.Encode(apperrors.BaseError{Error: err.Error()}); encodeErr != nil {
					fmt.Println(encodeErr.Error())
				}
				return
			}

			err = s.ReturnBook(r.Context(), loanID, body.ReturnDate)

			if err != nil {
				ProcessServiceErr(
					w,
					err,
					[]error{
						context.Canceled,
						context.DeadlineExceeded,
						apperrors.ErrNotFound,
					},
				)
				return
			}

			encoder := json.NewEncoder(w)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			if err = encoder.Encode(ReturnBookResponse{Success: true}); err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}
