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

func GetUserLoans(s LoanService) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				w.WriteHeader(apperrors.ErrDeadlineExceededCode)
			}
		default:
			userId, err := GetPathID(w, req.PathValue("userId"))

			if err != nil {
				return
			}

			userLoans, err := s.GetUserLoans(ctx, userId)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			encoder := json.NewEncoder(w)
			if err := encoder.Encode(dto.UserLoansModelToDTO(userLoans)); err != nil {
				fmt.Println(err)
			}
		}
	}
}
