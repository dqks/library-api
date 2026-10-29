package service

import "context"

func (s *loanService) ReturnBook(ctx context.Context, loanID int, returnDate *string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return s.repo.ReturnBook(ctx, loanID, returnDate)
	}
}
