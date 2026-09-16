package service

import (
	"context"
	"library-api/internal/model"
	"library-api/internal/repository"
	"library-api/internal/requests"
)

func GetBooks(ctx context.Context, params requests.GetBooksQueryParams) []*model.Book {
	return repository.GetBooks(ctx, params)
}
