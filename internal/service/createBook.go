package service

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/repository"
	"library-api/internal/requests"
)

func CreateBook(ctx context.Context, req requests.CreateBookRequest) (*model.Book, error) {
	if req.Author == nil || req.Available == nil || req.Title == nil || req.Year == nil {
		return nil, apperrors.ErrRequiredFields
	}
	book := repository.CreateBook(ctx, req)
	return book, nil
}
