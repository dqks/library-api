package postgres

import (
	"context"
	"errors"
	"fmt"
	"library-api/internal/apperrors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type loan struct {
	ID         int         `db:"id"`
	BookID     int         `db:"book_id"`
	UserID     int         `db:"user_id"`
	LoanDate   pgtype.Date `db:"loan_date"`
	ReturnDate pgtype.Date `db:"return_date"`
}

func (r *BookRepository) ReturnBook(ctx context.Context, loanID int, returnDate *string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		err := pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
			var loan loan
			err := tx.QueryRow(
				ctx,
				`
				SELECT id, book_id, user_id, loan_date, return_date FROM loans
				WHERE id = $1 FOR UPDATE
				`,
				loanID,
			).Scan(
				&loan.ID,
				&loan.BookID,
				&loan.UserID,
				&loan.LoanDate,
				&loan.ReturnDate,
			)

			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return fmt.Errorf("%w", apperrors.ErrNotFound)
				}
				return err
			}

			if loan.ReturnDate.Valid {
				return fmt.Errorf("%w: книгу уже вернули", apperrors.ErrReturnBook)
			}

			var tag pgconn.CommandTag

			if returnDate != nil {
				tag, err = tx.Exec(
					ctx,
					`
					UPDATE loans SET return_date = $1 WHERE id = $2
					`,
					returnDate,
					loanID,
				)

			} else {
				tag, err = tx.Exec(
					ctx,
					`
					UPDATE loans SET return_date = CURRENT_DATE WHERE id = $1
					`,
					loanID,
				)
			}

			if err != nil {
				return err
			}

			if tag.RowsAffected() == 0 {
				return apperrors.ErrNotFound
			}

			tag, err = tx.Exec(
				ctx,
				`
				UPDATE books SET available = true WHERE id = $1
				`,
				loan.BookID,
			)

			if err != nil {
				return err
			}

			if tag.RowsAffected() == 0 {
				return apperrors.ErrNotFound
			}

			return nil
		})
		if err != nil {
			return err
		}
		return nil
	}
}
