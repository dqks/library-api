package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"library-api/internal/apperrors"
	"library-api/internal/service"
	"net/http"
	"strconv"
)

func GetBookByID(s service.BookService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err())
		default:
			id, err := strconv.Atoi(r.PathValue("id"))

			if err != nil || id <= 0 {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(apperrors.ErrIncorrectPathValueCode)
				json.NewEncoder(w).Encode(apperrors.BaseError{Error: apperrors.ErrIncorrectPathValue.Error()})
				return
			}

			book, err := s.GetBookByID(ctx, id)

			if err != nil {
				code := apperrors.CheckErrors([]error{
					apperrors.ErrNotFound,
					context.DeadlineExceeded,
					context.Canceled,
				}, err)

				if code == -1 {
					return
				}

				w.Header().Add("Content-type", "application/json")
				w.WriteHeader(code)

				json.NewEncoder(w).Encode(apperrors.BaseError{Error: err.Error()})
				return
			}

			w.Header().Add("Content-Type", "application/json")
			encoder := json.NewEncoder(w)
			err = encoder.Encode(BookModelToDTO(&book))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}
}
