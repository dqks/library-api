package postgres

import (
	"context"
	"fmt"
	"library-api/internal/apperrors"

	"github.com/jackc/pgx/v5"
)

func (r *BookRepository) LoanBook(ctx context.Context, bookID, userID int) (int, error) {
	select {
	case <-ctx.Done():
		return -1, ctx.Err()
	default:
		id := -1

		err := pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
			// PostgreSQL защищает UPDATE операции он concurrency
			tag, err := tx.Exec(
				ctx,
				`
				UPDATE books SET available = false
				WHERE available = true AND id = $1
				`,
				bookID,
			)

			if err != nil {
				return err
			}

			if tag.RowsAffected() == 0 {
				return apperrors.ErrBookUnavailable
			}

			err = tx.QueryRow(
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
				return fmt.Errorf("%w: %v", apperrors.ErrInternal, err)
			}

			if id == -1 {
				return apperrors.ErrInternal
			}

			return nil
		})

		if err != nil {
			return -1, err
		}

		return id, nil
	}
}
