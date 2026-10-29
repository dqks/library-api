package postgres

import (
	"context"
	"errors"
	"fmt"
	"library-api/internal/apperrors"

	"github.com/jackc/pgx/v5"
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

			fmt.Println(loan)
			fmt.Println(loan.ReturnDate.Valid)

			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return fmt.Errorf("%w", apperrors.ErrNotFound)
				}
				return err
			}

			// var date string = *returnDate
			// if returnDate == nil {
			// 	date = "CURRENT_DATE"
			// }
			return nil
		})
		return err
	}
}
