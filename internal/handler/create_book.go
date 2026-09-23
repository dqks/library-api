package handler

import (
	"context"
	"encoding/json"
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

func CreateBook(s service.BookService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err().Error())
		default:
			var req CreateBookBody
			decoder := json.NewDecoder(r.Body)
			decoder.DisallowUnknownFields()

			if err := decoder.Decode(&req); err != nil {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(apperrors.ErrInvalidValuesCode)
				encoder := json.NewEncoder(w)
				if err := encoder.Encode(apperrors.BaseError{Error: err.Error()}); err != nil {
					fmt.Println(apperrors.ErrInternal.Error())
				}
				return
			}

			book, err := s.CreateBook(ctx, service.CreateBookInput{
				Title:     req.Title,
				Author:    req.Author,
				Year:      req.Year,
				Available: req.Available,
			})

			if err != nil {
				code := apperrors.CheckErrors([]error{
					context.DeadlineExceeded,
					context.Canceled,
					apperrors.ErrRequiredFields,
					apperrors.ErrInvalidValues,
				}, err)

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
			if err := encoder.Encode(BookModelToDTO(&book)); err != nil {
				fmt.Println(apperrors.ErrInternal.Error())
			}
		}
	}
}
