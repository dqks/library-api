package dto

import (
	"library-api/internal/model"
	"time"
)

type innerLoan struct {
	LoanID     int     `json:"loan_id"`
	LoanDate   string  `json:"loan_date"`
	DueDate    string  `json:"due_date"`
	ReturnDate *string `json:"return_date"`
	Book       BookDTO `json:"book"`
}

type UserLoanDTO struct {
	UserID   int         `json:"user_id"`
	Username string      `json:"user_name"`
	Loans    []innerLoan `json:"loans"`
}

func UserLoansModelToDTO(ul []model.UserLoan) UserLoanDTO {
	if ul == nil {
		return UserLoanDTO{}
	}

	var ulDTO UserLoanDTO
	ulDTO.UserID = ul[0].UserID
	ulDTO.Username = ul[0].Username

	for i := range ul {
		var loan innerLoan
		if ul[i].ReturnDate == nil {
			loan = innerLoan{
				LoanID:     ul[i].LoanID,
				LoanDate:   ul[i].LoanDate.Format(time.DateOnly),
				DueDate:    ul[i].DueDate.Format(time.DateOnly),
				ReturnDate: nil,
				Book: BookDTO{
					ID:        ul[i].Book.ID,
					Title:     ul[i].Book.Title,
					Author:    ul[i].Book.Author,
					Year:      ul[i].Book.Year,
					Available: ul[i].Book.Available,
				},
			}
		} else {
			returnDate := ul[i].ReturnDate.Format(time.DateOnly)
			loan = innerLoan{
				LoanID:     ul[i].LoanID,
				LoanDate:   ul[i].LoanDate.Format(time.DateOnly),
				ReturnDate: &returnDate,
			}
		}
		ulDTO.Loans = append(ulDTO.Loans, loan)
	}
	return ulDTO
}
