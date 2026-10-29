package tests

import (
	"context"
	"library-api/internal/model"
	"library-api/internal/service"
)

type ServiceMock struct {
	book            model.Book
	user            model.User
	books           []model.Book
	err             error
	callCount       int
	id              int
	bookID          int
	userID          int
	createBookInput service.CreateBookInput
	createUserInput service.CreateUserInput
	editInput       service.EditBookInput
	getParams       service.GetBooksQueryParams
	ctx             context.Context
}

func (s *ServiceMock) CreateBook(ctx context.Context, req service.CreateBookInput) (model.Book, error) {
	s.callCount++
	s.ctx = ctx
	s.createBookInput = req
	return s.book, s.err
}

func (s *ServiceMock) DeleteBookByID(ctx context.Context, id int) error {
	s.callCount++
	s.ctx = ctx
	s.id = id
	return s.err
}

func (s *ServiceMock) EditBookByID(ctx context.Context, id int, req service.EditBookInput) (model.Book, error) {
	s.callCount++
	s.ctx = ctx
	s.editInput = req
	s.id = id
	return s.book, s.err
}

func (s *ServiceMock) GetBookByID(ctx context.Context, id int) (model.Book, error) {
	s.callCount++
	s.ctx = ctx
	s.id = id
	return s.book, s.err
}

func (s *ServiceMock) GetBooks(ctx context.Context, params service.GetBooksQueryParams) ([]model.Book, error) {
	s.callCount++
	s.ctx = ctx
	s.getParams = params
	return s.books, s.err
}

func (s *ServiceMock) LoanBook(ctx context.Context, bookID, userID int) error {
	s.callCount++
	s.ctx = ctx
	s.userID = userID
	s.bookID = bookID
	return s.err
}

func (s *ServiceMock) CreateUser(ctx context.Context, input service.CreateUserInput) (model.User, error) {
	s.callCount++
	s.ctx = ctx
	s.createUserInput = input
	return s.user, s.err
}

func (s *ServiceMock) GetUser(ctx context.Context, id int) (model.User, error) {
	s.callCount++
	s.ctx = ctx
	s.id = id
	return s.user, s.err
}
