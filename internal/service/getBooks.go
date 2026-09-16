package service

import (
	"context"
	"library-api/internal/model"
	"library-api/internal/repository"
)

func GetBooks(ctx context.Context, params repository.GetBooksQueryParams) []*model.Book {
	return repository.GetBooks(ctx, params)
}
