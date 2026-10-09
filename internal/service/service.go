package service

import (
	"context"
	"library-api/internal/model"
	"library-api/internal/repository"
)

type Service interface {
	CreateBook(ctx context.Context, req CreateBookInput) (model.Book, error)
	DeleteBookByID(ctx context.Context, id int) error
	EditBookByID(ctx context.Context, id int, req EditBookInput) (model.Book, error)
	GetBookByID(ctx context.Context, id int) (model.Book, error)
	GetBooks(ctx context.Context, params GetBooksQueryParams) ([]model.Book, error)
}

type Repository interface {
	CreateBook(ctx context.Context, payload repository.CreateBookPayload) (model.Book, error)
	DeleteBookByID(ctx context.Context, id int) error
	EditBookByID(ctx context.Context, id int, p repository.EditBookPayload) (model.Book, error)
	GetBookByID(ctx context.Context, id int) (model.Book, error)
	GetBooks(ctx context.Context, params repository.GetBooksPayload) ([]model.Book, error)
}

type BookService struct {
	repo Repository
}

func Create(repo Repository) Service {
	if repo == nil {
		return nil
	}

	return &BookService{repo: repo}
}
