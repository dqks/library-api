package handler

import (
	"encoding/json"
	"fmt"
	"library-api/internal/apperrors"
	"library-api/internal/requests"
	"library-api/internal/service"
	"net/http"
)

func CreateBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req requests.CreateBookRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(400)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(apperrors.BaseError{Error: err.Error()}); err != nil {
			fmt.Println(apperrors.ErrInternal.Error())
		}
		return
	}

	book, err := service.CreateBook(ctx, req)

	if err != nil {
		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(400)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(apperrors.BaseError{Error: err.Error()}); err != nil {
			fmt.Println(apperrors.ErrInternal.Error())
		}
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(201)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(book.ToDTO()); err != nil {
		fmt.Println(apperrors.ErrInternal.Error())
	}
}
