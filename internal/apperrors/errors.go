package apperrors

import "errors"

var (
	ErrRequiredFields = errors.New("Отсутствуют обязательные поля")
	ErrNotFound       = errors.New("Не найдено")
)

type BaseError struct {
	Error string `json:"error"`
}
