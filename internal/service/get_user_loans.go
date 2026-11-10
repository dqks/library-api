package service

import (
	"context"
	"library-api/internal/model"
)

func (s *loanService) GetUserLoans(ctx context.Context, userID int) ([]model.UserLoan, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.GetUserLoans(ctx, userID)
	}
}
