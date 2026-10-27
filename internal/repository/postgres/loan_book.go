package postgres

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/repository"
)

func (r *BookRepository) LoanBook(ctx context.Context, bookID, userID int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, err := r.EditBookByID(ctx, bookID, repository.EditBookPayload{
			Available: new(false),
		})

		if err != nil {
			return err
		}

		tag, err := r.pool.Exec(
			ctx,
			`
			INSERT INTO loans  (book_id, user_id)
			VALUES ($1, $2)
			RETURNING book_id, user_id, loan_date
			`,
			bookID,
			userID,
		)

		if err != nil {
			return err
		}

		if tag.RowsAffected() == 0 {
			return apperrors.ErrNotFound
		}

		return nil
	}
}
