package dto

import (
	"library-api/internal/model"
	"time"
)

type LoanDTO struct {
	ID         int     `json:"id"`
	BookID     int     `json:"book_id"`
	UserID     int     `json:"user_id"`
	LoanDate   string  `json:"loan_date"`
	ReturnDate *string `json:"return_date"`
}

func LoanModelToDTO(loan *model.Loan) *LoanDTO {
	emptyTime := time.Time{}
	if (*loan.ReturnDate).Equal(emptyTime) {
		return &LoanDTO{
			ID:         loan.ID,
			BookID:     loan.BookID,
			UserID:     loan.UserID,
			LoanDate:   loan.LoanDate.Format(time.DateOnly),
			ReturnDate: nil,
		}
	} else {
		returnDate := loan.ReturnDate.Format(time.DateOnly)
		return &LoanDTO{
			ID:         loan.ID,
			BookID:     loan.BookID,
			UserID:     loan.UserID,
			LoanDate:   loan.LoanDate.Format(time.DateOnly),
			ReturnDate: &returnDate,
		}
	}
}
