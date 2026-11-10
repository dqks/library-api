package postgres

import (
	"context"
	"errors"
	"fmt"
	"library-api/internal/apperrors"
	"library-api/internal/model"

	"github.com/jackc/pgx/v5"
)

func (r *BookRepository) GetUserLoans(ctx context.Context, userID int) ([]model.UserLoan, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		tag, err := r.pool.Exec(
			ctx,
			`
			SELECT id FROM users WHERE id = $1
			`,
			userID,
		)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, apperrors.ErrNotFound
			}
			return nil, err
		}

		if tag.RowsAffected() == 0 {
			return nil, apperrors.ErrNotFound
		}

		rows, err := r.pool.Query(
			ctx,
			`
			SELECT 
				loans.id, 
				users.id, 
				users.name, 
				loan_date, 
				return_date, 
				books.id, 
				books.title,
				books.author,
				books.year, 
				books.available,
				loans.due_date
			FROM loans
			INNER JOIN users ON users.id = loans.user_id
			INNER JOIN books ON books.id = loans.book_id
			WHERE users.id = $1;
			`,
			userID,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		userLoans := make([]model.UserLoan, 0)

		for rows.Next() {
			var loan model.UserLoan
			if err := rows.Scan(
				&loan.LoanID,
				&loan.UserID,
				&loan.Username,
				&loan.LoanDate,
				&loan.ReturnDate,
				&loan.Book.ID,
				&loan.Book.Title,
				&loan.Book.Author,
				&loan.Book.Year,
				&loan.Book.Available,
				&loan.DueDate,
			); err != nil {
				return nil, err
			}
			userLoans = append(userLoans, loan)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}

		if len(userLoans) == 0 {
			return nil, fmt.Errorf("%w: взятые книги", apperrors.ErrNotFound)
		}

		return userLoans, nil
	}
}
