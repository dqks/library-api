package model

import (
	"time"
)

type UserLoan struct {
	UserID     int
	Username   string
	LoanID     int
	LoanDate   time.Time
	ReturnDate *time.Time
	DueDate    time.Time
	Book       Book
}
