package apperrors

import "errors"

var (
	ErrRequiredFields = errors.New("отсутствуют обязательные поля")
	ErrUnexpected     = errors.New("произошла непредвиденная ошибка")
	ErrNotFound       = errors.New("не найдено")
)

type BaseError struct {
	Error string `json:"error"`
}
