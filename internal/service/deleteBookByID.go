package service

import (
	"context"
	"library-api/internal/repository"
)

func DeleteBookByID(ctx context.Context, id int) error {
	return repository.DeleteBookByID(ctx, id)
}
