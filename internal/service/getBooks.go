package service

import (
	"context"
	"library-api/internal/model"
	"library-api/internal/repository"
)

type GetBooksQueryParams struct {
	Available *bool
}

func GetBooks(ctx context.Context, params GetBooksQueryParams) ([]model.Book, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return repository.GetBooks(ctx, repository.GetBooksPayload{Available: params.Available})
	}
}
