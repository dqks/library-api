package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/service"
	"net/http"
	"strconv"
)

func GetBooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	default:
		availParam := false
		encoder := json.NewEncoder(w)
		var books []model.Book
		var err error
		if availQuery := r.URL.Query().Get("available"); availQuery != "" {
			availParam, err = strconv.ParseBool(availQuery)
			if err != nil {
				w.WriteHeader(400)
				err := encoder.Encode(apperrors.BaseError{Error: err.Error()})
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				return
			}
			books, err = service.GetBooks(ctx, service.GetBooksQueryParams{Available: &availParam})
		} else {
			books, err = service.GetBooks(ctx, service.GetBooksQueryParams{Available: nil})
		}

		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			} else if errors.Is(err, context.DeadlineExceeded) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(504)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(500)
			}
			encoder.Encode(apperrors.BaseError{Error: err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = encoder.Encode(model.BookDomainListToDTO(books))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
