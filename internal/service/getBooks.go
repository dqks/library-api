package service

import (
	"context"
	"library-api/internal/model"
	"library-api/internal/repository"
)

func GetBooks(ctx context.Context) []*model.Book {
	return repository.GetBooks(ctx)
}
