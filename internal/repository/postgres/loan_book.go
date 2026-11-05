package postgres

import (
	"context"
	"fmt"
	"library-api/internal/apperrors"
	"library-api/internal/model"

	"github.com/jackc/pgx/v5"
)

func (r *BookRepository) LoanBook(ctx context.Context, bookID, userID int) (model.Loan, error) {
	select {
	case <-ctx.Done():
		return model.Loan{}, ctx.Err()
	default:
		var loan model.Loan
		err := pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
			id := -1
			err := tx.QueryRow(
				ctx,
				`
				SELECT id FROM books WHERE id = $1 FOR UPDATE;
				`,
				bookID,
			).Scan(
				&id,
			)
			if id == -1 {
				return apperrors.ErrNotFound
			}

			userLoanCount := 0
			err = tx.QueryRow(
				ctx,
				`
				SELECT COUNT(*) FROM loans
				WHERE user_id = $1 AND return_date IS NULL
				GROUP BY user_id;
				`,
				userID,
			).Scan(
				&userLoanCount,
			)
			if userLoanCount >= 5 {
				return apperrors.ErrUserLoanLimit
			}

			// PostgreSQL защищает UPDATE операции от concurrency
			tag, err := tx.Exec(
				ctx,
				`
				UPDATE books SET available = false
				WHERE available = true AND id = $1;
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
				RETURNING id, loan_date, return_date
				`,
				bookID,
				userID,
			).Scan(
				&loan.ID,
				&loan.LoanDate,
				&loan.ReturnDate,
			)

			if err != nil {
				return fmt.Errorf("%w: %v", apperrors.ErrInternal, err)
			}

			if loan.ID == 0 {
				return apperrors.ErrInternal
			}

			return nil
		})

		if err != nil {
			return model.Loan{}, err
		}

		return loan, nil
	}
}
