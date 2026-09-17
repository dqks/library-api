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

func EditBookByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		fmt.Println(apperrors.ErrContext.Error())
	default:
		id, err := strconv.Atoi(r.PathValue("id"))
		w.Header().Add("Content-Type", "application/json")

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(apperrors.BaseError{Error: apperrors.ErrIncorrectPathValue.Error()})
			return
		}

		var body service.EditBookRequest

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err = decoder.Decode(&body); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(apperrors.BaseError{Error: err.Error()})
			return
		}

		newBook, err := service.EditBookByID(ctx, id, body)

		if err != nil {
			if errors.Is(err, apperrors.ErrNotFound) {
				w.WriteHeader(404)
			} else if errors.Is(err, apperrors.ErrRequiredFields) || errors.Is(err, apperrors.ErrInvalidValues) {
				w.WriteHeader(400)
			}
			json.NewEncoder(w).Encode(apperrors.BaseError{Error: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(newBook.ToDTO())
	}
}
