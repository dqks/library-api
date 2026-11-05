package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"library-api/internal/apperrors"
	"library-api/internal/handler/dto"
	"net/http"
)

func GetLoan(s LoanService) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				w.WriteHeader(apperrors.ErrDeadlineExceededCode)
			}
			fmt.Println(ctx.Err())

		default:
			id, err := ValidatePathID(w, req.PathValue("id"))

			if err != nil {
				return
			}

			loan, err := s.GetLoan(ctx, id)

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

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			encoder := json.NewEncoder(w)
			if err = encoder.Encode(dto.LoanModelToDTO(&loan)); err != nil {
				fmt.Println(err)
			}
		}
	}
}
