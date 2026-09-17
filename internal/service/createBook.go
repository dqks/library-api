package service

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/repository"
	"strings"
)

type CreateBookRequest struct {
	Title     *string `json:"title"`
	Author    *string `json:"author"`
	Year      *uint16 `json:"year"`
	Available *bool   `json:"available"`
}

func CreateBook(ctx context.Context, req CreateBookRequest) (model.Book, error) {
	if req.Author == nil || req.Available == nil || req.Title == nil || req.Year == nil {
		return model.Book{}, apperrors.ErrRequiredFields
	}

	if len(strings.TrimSpace(*req.Author)) == 0 || len(strings.TrimSpace(*req.Title)) == 0 {
		return model.Book{}, apperrors.ErrInvalidValues
	}

	book := repository.CreateBook(ctx, repository.CreateBookPayload{
		Title:     req.Title,
		Author:    req.Author,
		Year:      req.Year,
		Available: req.Available,
	})
	return book, nil
}
