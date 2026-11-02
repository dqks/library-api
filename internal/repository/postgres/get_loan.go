package postgres

import (
	"context"
	"errors"
	"fmt"
	"library-api/internal/apperrors"
	"library-api/internal/model"

	"github.com/jackc/pgx/v5"
)

func (r *BookRepository) GetLoan(ctx context.Context, id int) (model.Loan, error) {
	select {
	case <-ctx.Done():
		return model.Loan{}, ctx.Err()
	default:
		var loan loan
		err := r.pool.QueryRow(
			ctx,
			`
			SELECT id, user_id, book_id, loan_date, return_date 
			FROM loans 
			WHERE id = $1
			`,
			id,
		).Scan(
			&loan.ID,
			&loan.UserID,
			&loan.BookID,
			&loan.LoanDate,
			&loan.ReturnDate,
		)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return model.Loan{}, fmt.Errorf("%w", apperrors.ErrNotFound)
			}
			return model.Loan{}, err
		}

		loanModel := model.Loan{
			ID:         loan.ID,
			BookID:     loan.BookID,
			UserID:     loan.UserID,
			LoanDate:   loan.LoanDate.Time,
			ReturnDate: &loan.ReturnDate.Time,
		}

		return loanModel, nil
	}
}
