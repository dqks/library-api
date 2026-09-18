package model

import (
	"library-api/internal/apperrors"
	"strings"
	"time"
)

type Book struct {
	ID        int
	Title     string
	Author    string
	Year      uint16
	Available bool
}

func ValidateBookFields(title *string, author *string, year *uint16) error {

	if (author != nil && len(strings.TrimSpace(*author)) == 0) || (title != nil && len(strings.TrimSpace(*title)) == 0) {
		return apperrors.ErrInvalidValues
	}

	if year != nil && *year > uint16(time.Now().Year()) {
		return apperrors.ErrInvalidValues
	}
	return nil
}
