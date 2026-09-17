package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"library-api/internal/apperrors"
	"library-api/internal/service"
	"net/http"
	"strconv"
)

type deleteBookResponse struct {
	Success bool `json:"success"`
}

func DeleteBookByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err().Error())
	default:
		id, err := strconv.Atoi(r.PathValue("id"))
		w.Header().Add("Content-type", "application/json")

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(
				apperrors.BaseError{Error: apperrors.ErrIncorrectPathValue.Error()},
			)
			return
		}

		err = service.DeleteBookByID(ctx, id)

		if err != nil {
			if errors.Is(err, apperrors.ErrNotFound) {
				w.WriteHeader(404)
			} else {
				w.WriteHeader(500)
			}
			json.NewEncoder(w).Encode(
				apperrors.BaseError{Error: err.Error()},
			)
			return
		}

		json.NewEncoder(w).Encode(
			deleteBookResponse{Success: true},
		)
	}

}
