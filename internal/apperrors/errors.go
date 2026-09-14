package apperrors

import "errors"

var (
	ErrRequiredFields = errors.New("Заказ не найден")
)
