package service

import (
	"context"
)

func (s *loanService) LoanBook(ctx context.Context, bookID, userID int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, err := s.repo.GetUser(ctx, userID)

		if err != nil {
			return err
		}

		_, err = s.repo.LoanBook(ctx, bookID, userID)

		return err
	}
}
