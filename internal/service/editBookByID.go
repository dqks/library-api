package service

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/repository"
	"library-api/internal/requests"
	"strings"
)

func EditBookByID(ctx context.Context, id int, req requests.EditBookRequest) (*model.Book, error) {
	if req.Author == nil && req.Available == nil && req.Title == nil && req.Year == nil {
		return nil, apperrors.ErrRequiredFields
	}

	if (req.Author != nil && len(strings.TrimSpace(*req.Author)) == 0) || (req.Author != nil && len(strings.TrimSpace(*req.Title)) == 0) {
		return nil, apperrors.ErrInvalidValues
	}

	return repository.EditBookByID(ctx, id, req)
}
