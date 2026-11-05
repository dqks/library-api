package dto

import (
	"fmt"
	"library-api/internal/model"
	"time"
)

type UserLoanDTO struct {
	LoanID     int     `json:"loan_id"`
	Username   string  `json:"user_name"`
	LoanDate   string  `json:"loan_date"`
	ReturnDate *string `json:"return_date"`
}

func UserLoanModelToDTO(ul *model.UserLoan) *UserLoanDTO {
	if ul.ReturnDate == nil {
		return &UserLoanDTO{
			LoanID:     ul.LoanID,
			Username:   ul.Username,
			LoanDate:   ul.LoanDate.Format(time.DateOnly),
			ReturnDate: nil,
		}
	}
	fmt.Println("before	returnDate := ul.ReturnDate.Format(time.DateOnly)")
	returnDate := ul.ReturnDate.Format(time.DateOnly)
	fmt.Println("after	returnDate := ul.ReturnDate.Format(time.DateOnly)")
	return &UserLoanDTO{
		LoanID:     ul.LoanID,
		Username:   ul.Username,
		LoanDate:   ul.LoanDate.Format(time.DateOnly),
		ReturnDate: &returnDate,
	}
}

func UserLoansModelToDTO(ul []model.UserLoan) []UserLoanDTO {
	if ul == nil {
		fmt.Println("ul == nil")
		return nil
	}

	ulDTO := make([]UserLoanDTO, 0, len(ul))
	for i := range ul {
		ulDTO = append(ulDTO, *UserLoanModelToDTO(&ul[i]))
	}

	return ulDTO
}
