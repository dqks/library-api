package tests

import (
	"context"
	"library-api/internal/model"
	"library-api/internal/service"
)

type ServiceMock struct {
	book           model.Book
	books          []model.Book
	err            error
	gotIntoService bool
	id             int
	createInput    service.CreateBookInput
	editInput      service.EditBookInput
	getParams      service.GetBooksQueryParams
}

func (s *ServiceMock) CreateBook(ctx context.Context, req service.CreateBookInput) (model.Book, error) {
	s.gotIntoService = true
	s.createInput = req
	return s.book, s.err
}

func (s *ServiceMock) DeleteBookByID(ctx context.Context, id int) error {
	s.gotIntoService = true
	s.id = id
	return s.err
}

func (s *ServiceMock) EditBookByID(ctx context.Context, id int, req service.EditBookInput) (model.Book, error) {
	s.gotIntoService = true
	s.editInput = req
	s.id = id
	return s.book, s.err
}

func (s *ServiceMock) GetBookByID(ctx context.Context, id int) (model.Book, error) {
	s.gotIntoService = true
	s.id = id
	return s.book, s.err
}

func (s *ServiceMock) GetBooks(ctx context.Context, params service.GetBooksQueryParams) ([]model.Book, error) {
	s.gotIntoService = true
	s.getParams = params
	return s.books, s.err
}
