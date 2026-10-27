package service

import (
	"context"
	"library-api/internal/apperrors"
)

func (s *loanService) LoanBook(ctx context.Context, bookID, userID int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		book, err := s.repo.GetBookByID(ctx, bookID)

		if err != nil {
			return err
		}

		if !book.Available {
			return apperrors.ErrBookUnavailable
		}

		_, err = s.repo.GetUser(ctx, userID)

		if err != nil {
			return err
		}

		_, err = s.repo.LoanBook(ctx, bookID, userID)

		return err
	}
}
