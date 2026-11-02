package service

import (
	"context"
	"library-api/internal/model"
)

func (s *loanService) GetLoan(ctx context.Context, id int) (model.Loan, error) {
	select {
	case <-ctx.Done():
		return model.Loan{}, ctx.Err()
	default:
		return s.repo.GetLoan(ctx, id)
	}
}
