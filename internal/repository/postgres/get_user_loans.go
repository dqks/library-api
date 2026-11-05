package postgres

import (
	"context"
	"fmt"
	"library-api/internal/model"
)

func (r *BookRepository) GetUserLoans(ctx context.Context, userID int) ([]model.UserLoan, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		rows, err := r.pool.Query(
			ctx,
			`
			SELECT loans.id, users.name, loan_date, return_date
			FROM loans
			INNER JOIN users ON users.id = loans.user_id
			WHERE users.id = $1;
			`,
			userID,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		fmt.Println(rows.RawValues())

		userLoans := make([]model.UserLoan, 0)

		for rows.Next() {
			var loan model.UserLoan
			if err := rows.Scan(
				&loan.LoanID,
				&loan.Username,
				&loan.LoanDate,
				&loan.ReturnDate,
			); err != nil {
				return nil, err
			}
			userLoans = append(userLoans, loan)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}

		return userLoans, nil
	}
}
