package model

import "time"

type Loan struct {
	ID         int
	BookID     int
	UserID     int
	LoanDate   time.Time
	ReturnDate *time.Time
}
