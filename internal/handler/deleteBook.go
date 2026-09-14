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
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(
			apperrors.BaseError{Error: apperrors.ErrNotFound.Error()},
		)
		return
	}

	err = service.DeleteBookByID(ctx, id)

	if errors.Is(err, apperrors.ErrNotFound) {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(
			apperrors.BaseError{Error: apperrors.ErrNotFound.Error()},
		)
		return
	}

}
