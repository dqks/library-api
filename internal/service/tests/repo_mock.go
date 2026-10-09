package tests

import (
	"context"
	"library-api/internal/model"
	"library-api/internal/repository"
)

type MockRepo struct {
	book              model.Book
	err               error
	books             []model.Book
	id                int
	gotInside         bool
	getBooksPayload   repository.GetBooksPayload
	editBookPayload   repository.EditBookPayload
	createBookPayload repository.CreateBookPayload
}

func (r *MockRepo) CreateBook(ctx context.Context, p repository.CreateBookPayload) (model.Book, error) {
	r.gotInside = true
	r.createBookPayload = p
	return r.book, r.err
}

func (r *MockRepo) DeleteBookByID(ctx context.Context, id int) error {
	r.gotInside = true
	r.id = id
	return r.err
}

func (r *MockRepo) EditBookByID(ctx context.Context, id int, p repository.EditBookPayload) (model.Book, error) {
	r.gotInside = true
	r.id = id
	r.editBookPayload = p
	return r.book, r.err
}

func (r *MockRepo) GetBookByID(ctx context.Context, id int) (model.Book, error) {
	r.gotInside = true
	r.id = id
	return r.book, r.err
}

func (r *MockRepo) GetBooks(ctx context.Context, p repository.GetBooksPayload) ([]model.Book, error) {
	r.gotInside = true
	r.getBooksPayload = p
	return r.books, r.err
}
