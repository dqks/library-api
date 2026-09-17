package service

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/repository"
)

func GetBookByID(ctx context.Context, id int) (model.Book, error) {
	select {
	case <-ctx.Done():
		return model.Book{}, apperrors.ErrContext
	default:
		return repository.GetBookByID(ctx, id)
	}
}
