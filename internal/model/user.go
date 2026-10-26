package model

import (
	"library-api/internal/apperrors"
	"strings"
)

type User struct {
	ID   int
	Name string
}

func ValidateUserFields(name *string) error {
	if name == nil || len(strings.TrimSpace(*name)) == 0 {
		return apperrors.ErrInvalidValues
	}
	return nil
}
