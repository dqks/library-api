package service

import (
	"context"
	"library-api/internal/model"
	repository "library-api/internal/repository/inmemory"
)

type Repository interface {
	CreateBook(ctx context.Context, payload repository.CreateBookPayload) (model.Book, error)
	DeleteBookByID(ctx context.Context, id int) error
	EditBookByID(ctx context.Context, id int, p repository.EditBookPayload) (model.Book, error)
	GetBookByID(ctx context.Context, id int) (model.Book, error)
	GetBooks(ctx context.Context, params repository.GetBooksPayload) ([]model.Book, error)
}

type bookService struct {
	repo Repository
}

func Create(repo Repository) *bookService {
	if repo == nil {
		return nil
	}

	return &bookService{repo: repo}
}
