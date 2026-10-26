package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"library-api/internal/apperrors"
	"library-api/internal/handler/dto"
	"library-api/internal/service"
	"net/http"
)

type CreateUserBody struct {
	Name *string `json:"name"`
}

func CreateUser(s UserService) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				w.WriteHeader(apperrors.ErrDeadlineExceededCode)
			}
		default:
			var body CreateUserBody

			err := DecodeBody(r.Body, &body, w)

			if err != nil {
				return
			}

			user, err := s.CreateUser(ctx, service.CreateUserInput{
				Name: body.Name,
			})

			if err != nil {
				ProcessServiceErr(w, err)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(201)
			encoder := json.NewEncoder(w)
			if err = encoder.Encode(dto.UserModelToDTO(&user)); err != nil {
				fmt.Println(apperrors.ErrInternal.Error())
			}
		}
	}
}
