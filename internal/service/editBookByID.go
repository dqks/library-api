package service

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/repository"
)

type EditBookRequest struct {
	Title     *string `json:"title"`
	Author    *string `json:"author"`
	Year      *uint16 `json:"year"`
	Available *bool   `json:"available"`
}

func EditBookByID(ctx context.Context, id int, req EditBookRequest) (model.Book, error) {
	select {
	case <-ctx.Done():
		return model.Book{}, ctx.Err()
	default:
		if req.Author == nil && req.Available == nil && req.Title == nil && req.Year == nil {
			return model.Book{}, apperrors.ErrRequiredFields
		}

		if err := model.ValidateBookFields(req.Title, req.Author, req.Year); err != nil {
			return model.Book{}, err
		}

		return repository.EditBookByID(ctx, id, repository.EditBookPayload{
			Title:     req.Title,
			Author:    req.Author,
			Year:      req.Year,
			Available: req.Available,
		})
	}
}
