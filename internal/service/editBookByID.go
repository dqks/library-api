package service

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/repository"
	"strings"
)

type EditBookRequest struct {
	Title     *string `json:"title"`
	Author    *string `json:"author"`
	Year      *uint16 `json:"year"`
	Available *bool   `json:"available"`
}

func EditBookByID(ctx context.Context, id int, req EditBookRequest) (*model.Book, error) {
	if req.Author == nil && req.Available == nil && req.Title == nil && req.Year == nil {
		return nil, apperrors.ErrRequiredFields
	}

	if (req.Author != nil && len(strings.TrimSpace(*req.Author)) == 0) || (req.Title != nil && len(strings.TrimSpace(*req.Title)) == 0) {
		return nil, apperrors.ErrInvalidValues
	}

	return repository.EditBookByID(ctx, id, repository.EditBookPayload{
		Title:     req.Title,
		Author:    req.Author,
		Year:      req.Year,
		Available: req.Available,
	})
}
