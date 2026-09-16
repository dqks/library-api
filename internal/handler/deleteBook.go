package handler

import (
	"encoding/json"
	"library-api/internal/apperrors"
	"library-api/internal/service"
	"net/http"
	"strconv"
)

func DeleteBookByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		w.Header().Add("Content-type", "application-json")
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(
			apperrors.BaseError{Error: apperrors.ErrNotFound.Error()},
		)
		return
	}

	err = service.DeleteBookByID(ctx, id)

	if err != nil {
		w.Header().Add("Content-type", "application-json")
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(
			apperrors.BaseError{Error: err.Error()},
		)
		return
	}

}
