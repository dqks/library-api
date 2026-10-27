package service

import (
	"context"
	"library-api/internal/model"
)

type LoanRepository interface {
	GetBookByID(ctx context.Context, id int) (model.Book, error)
	GetUser(ctx context.Context, id int) (model.User, error)
	LoanBook(ctx context.Context, bookID, userID int) error
}

type loanService struct {
	repo LoanRepository
}

func CreateLoanService(repo LoanRepository) *loanService {
	if repo == nil {
		return nil
	}

	return &loanService{repo: repo}
}
