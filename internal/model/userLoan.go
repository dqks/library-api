package model

import "time"

type UserLoan struct {
	LoanID     int
	Username   string
	LoanDate   time.Time
	ReturnDate *time.Time
}
