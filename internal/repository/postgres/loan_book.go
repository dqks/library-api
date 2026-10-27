package postgres

import (
	"context"
	"fmt"
	"library-api/internal/apperrors"
	"library-api/internal/repository"
)

func (r *BookRepository) LoanBook(ctx context.Context, bookID, userID int) (int, error) {
	select {
	case <-ctx.Done():
		return -1, ctx.Err()
	default:
		_, err := r.EditBookByID(ctx, bookID, repository.EditBookPayload{
			Available: new(false),
		})

		if err != nil {
			return -1, err
		}

		id := -1

		err = r.pool.QueryRow(
			ctx,
			`
			INSERT INTO loans (book_id, user_id)
			VALUES ($1, $2)
			RETURNING id
			`,
			bookID,
			userID,
		).Scan(
			&id,
		)

		if err != nil {
			return -1, fmt.Errorf("%w: %v", apperrors.ErrInternal, err)
		}

		if id == -1 {
			return -1, apperrors.ErrInternal
		}

		return id, nil
	}
}
