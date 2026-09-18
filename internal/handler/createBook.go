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

func CreateBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err().Error())
	default:
		var req service.CreateBookRequest
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
			if errors.Is(err, context.DeadlineExceeded) {
				w.WriteHeader(504)
				return
			} else if errors.Is(err, context.Canceled) {
				return
			} else if errors.Is(err, apperrors.ErrRequiredFields) || errors.Is(err, apperrors.ErrInvalidValues) {
				w.Header().Add("Content-type", "application/json")
				w.WriteHeader(400)
				encoder := json.NewEncoder(w)
				if err := encoder.Encode(apperrors.BaseError{Error: err.Error()}); err != nil {
					fmt.Println(apperrors.ErrInternal.Error())
				}
				return
			} else {
				w.Header().Add("Content-type", "application/json")
				w.WriteHeader(500)
				encoder := json.NewEncoder(w)
				if err := encoder.Encode(apperrors.BaseError{Error: err.Error()}); err != nil {
					fmt.Println(apperrors.ErrInternal.Error())
				}
				return
			}
		}

		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(201)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(book.ToDTO()); err != nil {
			fmt.Println(apperrors.ErrInternal.Error())
		}
	}
}
