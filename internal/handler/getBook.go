package handler

import (
	"encoding/json"
	"errors"
	"library-api/internal/apperrors"
	"library-api/internal/service"
	"net/http"
	"strconv"
)

func GetBookByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(r.PathValue("id"))
	w.Header().Add("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(apperrors.BaseError{Error: apperrors.ErrIncorrectPathValue.Error()})
		return
	}

	book, err := service.GetBookByID(ctx, id)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			w.WriteHeader(404)
		} else {
			w.WriteHeader(500)
		}
		json.NewEncoder(w).Encode(apperrors.BaseError{Error: err.Error()})
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(book.ToDTO())
}
