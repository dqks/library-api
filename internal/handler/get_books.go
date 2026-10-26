package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"library-api/internal/apperrors"
	"library-api/internal/handler/dto"
	"library-api/internal/model"
	"library-api/internal/service"
	"net/http"
	"strconv"
)

func GetBooks(s BookService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				w.WriteHeader(apperrors.ErrDeadlineExceededCode)
			}
			fmt.Println(ctx.Err())
		default:
			availParam := false
			encoder := json.NewEncoder(w)
			var books []model.Book
			var err error
			if availQuery := r.URL.Query().Get("available"); availQuery != "" {
				availParam, err = strconv.ParseBool(availQuery)
				if err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(apperrors.ErrIncorrectPathValueCode)
					err := encoder.Encode(apperrors.BaseError{Error: err.Error()})
					if err != nil {
						return
					}
					return
				}
				books, err = s.GetBooks(ctx, service.GetBooksQueryParams{Available: &availParam})
			} else {
				books, err = s.GetBooks(ctx, service.GetBooksQueryParams{Available: nil})
			}

			if err != nil {
				code := apperrors.CheckErrors([]error{
					context.DeadlineExceeded,
					context.Canceled,
				}, err)

				if code == -1 {
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(code)

				encoder.Encode(apperrors.BaseError{Error: err.Error()})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			err = encoder.Encode(dto.BookDomainListToDTO(books))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}
}
