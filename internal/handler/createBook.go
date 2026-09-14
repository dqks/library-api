package handler

import (
	"encoding/json"
	"library-api/internal/repository"
	"library-api/internal/service"
	"net/http"
)

type createBookError struct {
	ErrorMessage string `json:"errorMessage"`
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req repository.CreateBookRequest
	json.NewDecoder(r.Body).Decode(&req)

	err := service.CreateBook(ctx, req)

	w.Header().Add("Content-Type", "application-json")

	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(createBookError{ErrorMessage: err.Error()})
		return
	}

	w.WriteHeader(201)
}
