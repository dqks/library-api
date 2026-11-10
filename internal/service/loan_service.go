package service

import (
	"context"
	"library-api/internal/model"
)

type LoanRepository interface {
	GetUser(ctx context.Context, id int) (model.User, error)
	LoanBook(ctx context.Context, bookID, userID int) (model.Loan, error)
	ReturnBook(ctx context.Context, loanID int, returnDate *string) error
	GetLoan(ctx context.Context, id int) (model.Loan, error)
	GetUserLoans(ctx context.Context, userID int) ([]model.UserLoan, error)
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
