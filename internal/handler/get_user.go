package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"library-api/internal/apperrors"
	"net/http"
	"strconv"
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
			id, err := strconv.Atoi(r.PathValue("id"))

			if err != nil || id <= 0 {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(apperrors.ErrIncorrectPathValueCode)
				encoder := json.NewEncoder(w)
				if err := encoder.Encode(apperrors.BaseError{Error: err.Error()}); err != nil {
					fmt.Println(err)
				}
			}

			user, err := s.GetUser(ctx, id)

			fmt.Println(user, id)
		}
	}
}
