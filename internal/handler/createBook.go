package handler

import (
	"encoding/json"
	"library-api/internal/apperrors"
	"library-api/internal/repository"
	"library-api/internal/service"
	"net/http"
)

func CreateBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Add("Content-Type", "application-json")

	var req repository.CreateBookRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(apperrors.BaseError{Error: err.Error()})
	}

	err := service.CreateBook(ctx, req)

	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(apperrors.BaseError{Error: err.Error()})
		return
	}

	w.WriteHeader(201)
}
