package handler

import (
	"encoding/json"
	"errors"
	"library-api/internal/apperrors"
	"library-api/internal/service"
	"net/http"
	"strconv"
)

func DeleteBookByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(
			apperrors.BaseError{Error: apperrors.ErrIncorrectPathValue.Error()},
		)
		return
	}

	err = service.DeleteBookByID(ctx, id)

	if err != nil {
		w.Header().Add("Content-type", "application/json")
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
		struct {
			success bool
		}{success: true},
	)
}
