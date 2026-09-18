package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"library-api/internal/apperrors"
	"library-api/internal/service"
	"net/http"
)

type CreateBookBody struct {
	Title     *string `json:"title"`
	Author    *string `json:"author"`
	Year      *uint16 `json:"year"`
	Available *bool   `json:"available"`
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err().Error())
	default:
		var req CreateBookBody
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

		book, err := service.CreateBook(ctx, service.CreateBookInput{
			Title:     req.Title,
			Author:    req.Author,
			Year:      req.Year,
			Available: req.Available,
		})

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				w.WriteHeader(504)
			} else if errors.Is(err, context.Canceled) {
				return
			} else if errors.Is(err, apperrors.ErrRequiredFields) || errors.Is(err, apperrors.ErrInvalidValues) {
				w.Header().Add("Content-type", "application/json")
				w.WriteHeader(400)
			} else {
				w.Header().Add("Content-type", "application/json")
				w.WriteHeader(500)
			}

			encoder := json.NewEncoder(w)
			if err := encoder.Encode(apperrors.BaseError{Error: err.Error()}); err != nil {
				fmt.Println(apperrors.ErrInternal.Error())
			}
			return
		}

		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(201)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(BookModelToDTO(&book)); err != nil {
			fmt.Println(apperrors.ErrInternal.Error())
		}
	}
}
