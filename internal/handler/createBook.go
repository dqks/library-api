package handler

import (
	"encoding/json"
	"fmt"
	"library-api/internal/apperrors"
	"library-api/internal/repository"
	"library-api/internal/service"
	"net/http"
)

func CreateBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req repository.CreateBookRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		w.WriteHeader(400)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(apperrors.BaseError{Error: err.Error()}); err != nil {
			fmt.Println(apperrors.ErrUnexpected.Error())
		}
		return
	}

	err := service.CreateBook(ctx, req)

	if err != nil {
		w.WriteHeader(400)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(apperrors.BaseError{Error: err.Error()}); err != nil {
			fmt.Println(apperrors.ErrUnexpected.Error())
		}
		return
	}

	w.WriteHeader(201)
}
