package service

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/repository"
	"library-api/internal/requests"
	"strings"
)

func CreateBook(ctx context.Context, req requests.CreateBookRequest) (*model.Book, error) {
	if req.Author == nil || req.Available == nil || req.Title == nil || req.Year == nil {
		return nil, apperrors.ErrRequiredFields
	}

	if len(strings.TrimSpace(*req.Author)) == 0 || len(strings.TrimSpace(*req.Title)) == 0 {
		return nil, apperrors.ErrInvalidValues
	}

	book := repository.CreateBook(ctx, req)
	return book, nil
}
