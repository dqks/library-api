package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"library-api/internal/apperrors"
	"library-api/internal/service"
	"net/http"
	"strconv"
)

func GetBookByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	default:
		id, err := strconv.Atoi(r.PathValue("id"))

		if err != nil || id <= 0 {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(apperrors.BaseError{Error: apperrors.ErrIncorrectPathValue.Error()})
			return
		}

		book, err := service.GetBookByID(ctx, id)

		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			} else if errors.Is(err, context.DeadlineExceeded) {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(504)
			} else if errors.Is(err, apperrors.ErrNotFound) {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(404)
			} else {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(500)
			}
			json.NewEncoder(w).Encode(apperrors.BaseError{Error: err.Error()})
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(book.ToDTO())
	}
}
