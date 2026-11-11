package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"library-api/internal/apperrors"
	"library-api/internal/handler/dto"
	"library-api/internal/service"
	"log/slog"
	"net/http"
)

type CreateBookBody struct {
	Title     *string `json:"title"`
	Author    *string `json:"author"`
	Year      *uint16 `json:"year"`
	Available *bool   `json:"available"`
}

func CreateBook(s BookService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				w.WriteHeader(apperrors.ErrDeadlineExceededCode)
			}
			fmt.Println(ctx.Err().Error())
		default:
			var body CreateBookBody
			decoder := json.NewDecoder(r.Body)
			decoder.DisallowUnknownFields()

			if err := decoder.Decode(&body); err != nil {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(apperrors.ErrInvalidValuesCode)
				encoder := json.NewEncoder(w)
				if err := encoder.Encode(apperrors.BaseError{Error: err.Error()}); err != nil {
					fmt.Println(apperrors.ErrInternal.Error())
				}
				return
			}

			book, err := s.CreateBook(ctx, service.CreateBookInput{
				Title:     body.Title,
				Author:    body.Author,
				Year:      body.Year,
				Available: body.Available,
			})

			if err != nil {
				code := apperrors.CheckErrors([]error{
					context.DeadlineExceeded,
					context.Canceled,
					apperrors.ErrRequiredFields,
					apperrors.ErrInvalidValues,
				}, err)

				slog.Error(
					"create book error",
					"body_title", *body.Title,
					"body_author", *body.Author,
					"body_year", *body.Year,
					"body_available", *body.Available,
					"err", err,
				)

				if code == -1 {
					return
				}

				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(code)

				encoder := json.NewEncoder(w)
				if err := encoder.Encode(apperrors.BaseError{Error: err.Error()}); err != nil {
					fmt.Println(apperrors.ErrInternal.Error())
				}
				return
			}

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(201)
			encoder := json.NewEncoder(w)
			if err := encoder.Encode(dto.BookModelToDTO(&book)); err != nil {
				fmt.Println(apperrors.ErrInternal.Error())
			}
		}
	}
}
