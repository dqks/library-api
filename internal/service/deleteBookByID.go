package service

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/repository"
)

func DeleteBookByID(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return apperrors.ErrContext
	default:
		return repository.DeleteBookByID(ctx, id)

	}
}
