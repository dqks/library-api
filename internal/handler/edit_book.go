package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"library-api/internal/apperrors"
	"library-api/internal/service"
	"net/http"
	"strconv"
)

type EditBookBody struct {
	Title     *string `json:"title"`
	Author    *string `json:"author"`
	Year      *uint16 `json:"year"`
	Available *bool   `json:"available"`
}

func EditBookByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err().Error())
	default:
		id, err := strconv.Atoi(r.PathValue("id"))

		if err != nil || id <= 0 {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(apperrors.BaseError{Error: apperrors.ErrIncorrectPathValue.Error()})
			return
		}

		var body EditBookBody

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err = decoder.Decode(&body); err != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(apperrors.BaseError{Error: err.Error()})
			return
		}

		newBook, err := service.EditBookByID(ctx, id, service.EditBookInput{
			Title:     body.Title,
			Author:    body.Author,
			Year:      body.Year,
			Available: body.Available,
		})

		if err != nil {
			if errors.Is(err, apperrors.ErrNotFound) {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(404)
			} else if errors.Is(err, apperrors.ErrRequiredFields) || errors.Is(err, apperrors.ErrInvalidValues) {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(400)
			} else if errors.Is(err, context.DeadlineExceeded) {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(504)
			} else if errors.Is(err, context.Canceled) {
				return
			} else {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(500)
			}
			json.NewEncoder(w).Encode(apperrors.BaseError{Error: err.Error()})
			return
		}

		w.Header().Add("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err = encoder.Encode(BookModelToDTO(&newBook))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
