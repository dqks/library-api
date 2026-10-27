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

func GetUser(s UserService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				w.WriteHeader(apperrors.ErrDeadlineExceededCode)
			}
		default:
			id, err := ValidatePathID(w, r.PathValue("id"))

			if err != nil {
				return
			}

			user, err := s.GetUser(ctx, id)

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
			w.WriteHeader(200)
			encoder := json.NewEncoder(w)
			if encoder.Encode(dto.UserModelToDTO(&user)); err != nil {
				fmt.Println(err)
			}
		}
	}
}
