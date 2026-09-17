package service

import (
	"context"
	"library-api/internal/repository"
)

func DeleteBookByID(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return repository.DeleteBookByID(ctx, id)
	}
}
