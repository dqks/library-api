package service

import (
	"context"
	"library-api/internal/model"
	"library-api/internal/repository"
)

type GetBooksQueryParams struct {
	Available bool
}

func GetBooks(ctx context.Context, params GetBooksQueryParams) []*model.Book {
	return repository.GetBooks(ctx, repository.GetBooksPayload{Available: params.Available})
}
