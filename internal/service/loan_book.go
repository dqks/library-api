package service

import (
	"context"
	"library-api/internal/model"
)

func (s *loanService) LoanBook(ctx context.Context, bookID, userID int) (model.Loan, error) {
	select {
	case <-ctx.Done():
		return model.Loan{}, ctx.Err()
	default:
		_, err := s.repo.GetUser(ctx, userID)

		if err != nil {
			return model.Loan{}, err
		}

		loan, err := s.repo.LoanBook(ctx, bookID, userID)

		return loan, err
	}
}
