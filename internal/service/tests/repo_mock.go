package tests

import (
	"context"
	"library-api/internal/model"
	repository "library-api/internal/repository"
)

type MockRepo struct {
	//return values
	book  model.Book
	err   error
	books []model.Book
	user  model.User
	//args
	bookID     int
	userID     int
	getBookID  int
	getUserID  int
	loanBookID int
	loanUserID int
	//logic
	gotInside bool
	// payloads
	getBooksPayload   repository.GetBooksPayload
	editBookPayload   repository.EditBookPayload
	createBookPayload repository.CreateBookPayload
	createUserPayload repository.CreateUserPayload
	//errs
	getUserErr  error
	loanBookErr error
}

func (r *MockRepo) CreateBook(ctx context.Context, p repository.CreateBookPayload) (model.Book, error) {
	r.gotInside = true
	r.createBookPayload = p
	return r.book, r.err
}

func (r *MockRepo) DeleteBookByID(ctx context.Context, id int) error {
	r.gotInside = true
	r.bookID = id
	return r.err
}

func (r *MockRepo) EditBookByID(ctx context.Context, id int, p repository.EditBookPayload) (model.Book, error) {
	r.gotInside = true
	r.bookID = id
	r.editBookPayload = p
	return r.book, r.err
}

func (r *MockRepo) GetBooks(ctx context.Context, p repository.GetBooksPayload) ([]model.Book, error) {
	r.gotInside = true
	r.getBooksPayload = p
	return r.books, r.err
}

func (r *MockRepo) CreateUser(ctx context.Context, p repository.CreateUserPayload) (model.User, error) {
	r.gotInside = true
	r.createUserPayload = p
	return r.user, r.err
}

func (r *MockRepo) GetBookByID(ctx context.Context, id int) (model.Book, error) {
	r.gotInside = true
	r.getBookID = id
	return r.book, r.err
}

func (r *MockRepo) GetUser(ctx context.Context, id int) (model.User, error) {
	r.gotInside = true
	r.getUserID = id
	return r.user, r.getUserErr
}

func (r *MockRepo) LoanBook(ctx context.Context, bookID, userID int) (int, error) {
	r.gotInside = true
	r.loanUserID = userID
	r.loanBookID = bookID
	return 1, r.loanBookErr
}
