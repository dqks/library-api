package service

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/repository"
	"library-api/internal/requests"
)

func EditBookByID(ctx context.Context, id int, req requests.EditBookRequest) (*model.Book, error) {
	if req.Author == nil && req.Available == nil && req.Title == nil && req.Year == nil {
		return nil, apperrors.ErrRequiredFields
	}
	return repository.EditBookByID(ctx, id, req)
}
