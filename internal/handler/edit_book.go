package handler

import (
	"context"
	"encoding/json"
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

func EditBookByID(s service.BookService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
				w.WriteHeader(apperrors.ErrIncorrectPathValueCode)
				json.NewEncoder(w).Encode(apperrors.BaseError{Error: err.Error()})
				return
			}

			newBook, err := s.EditBookByID(ctx, id, service.EditBookInput{
				Title:     body.Title,
				Author:    body.Author,
				Year:      body.Year,
				Available: body.Available,
			})

			if err != nil {
				code := apperrors.CheckErrors([]error{
					apperrors.ErrNotFound,
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
}
