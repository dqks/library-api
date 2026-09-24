package apperrors

import (
	"context"
	"errors"
)

var (
	ErrRequiredFields     = errors.New("отсутствуют обязательные поля")
	ErrInternal           = errors.New("произошла внутренняя ошибка")
	ErrNotFound           = errors.New("не найдено")
	ErrIncorrectPathValue = errors.New("неверное значение в пути")
	ErrInvalidValues      = errors.New("невалидные данные")
)

var (
	ErrRequiredFieldsCode     = 400
	ErrInternalCode           = 500
	ErrNotFoundCode           = 404
	ErrIncorrectPathValueCode = 400
	ErrInvalidValuesCode      = 400
	ErrDeadlineExceeded       = 504
	ErrCanceled               = -1
)

type BaseError struct {
	Error string `json:"error"`
}

var errorMap = map[error]int{
	ErrRequiredFields:        ErrRequiredFieldsCode,
	ErrInternal:              ErrInternalCode,
	ErrNotFound:              ErrNotFoundCode,
	ErrIncorrectPathValue:    ErrIncorrectPathValueCode,
	ErrInvalidValues:         ErrInvalidValuesCode,
	context.DeadlineExceeded: ErrDeadlineExceeded,
	context.Canceled:         ErrCanceled,
}

func CheckErrors(errorsToCheck []error, err error) int {
	errIndex := -1

	for i := range errorsToCheck {
		if errors.Is(errorsToCheck[i], err) {
			errIndex = i
		}
	}

	if errIndex == -1 {
		return ErrInternalCode
	}

	if errorMap[errorsToCheck[errIndex]] != 0 {
		return errorMap[errorsToCheck[errIndex]]
	}

	return ErrInternalCode
}
