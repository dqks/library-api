package service

import (
	"context"
	"library-api/internal/model"
	"library-api/internal/repository"
)

func GetBookByID(ctx context.Context, id int) (model.Book, error) {
	return repository.GetBookByID(ctx, id)
}
