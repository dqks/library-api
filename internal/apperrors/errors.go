package apperrors

import "errors"

var (
	ErrRequiredFields     = errors.New("отсутствуют обязательные поля")
	ErrInternal           = errors.New("произошла внутренняя ошибка")
	ErrNotFound           = errors.New("не найдено")
	ErrIncorrectPathValue = errors.New("неверное значение в пути")
)

type BaseError struct {
	Error string `json:"error"`
}
