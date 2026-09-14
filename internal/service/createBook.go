package service

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/repository"
)

func CreateBook(ctx context.Context, req repository.CreateBookRequest) error {
	if req.Author == nil || req.Available == nil || req.Title == nil || req.Year == nil {
		return apperrors.ErrRequiredFields
	}
	repository.CreateBook(ctx, req)
	return nil
}
