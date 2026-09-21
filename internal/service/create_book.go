package service

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/repository"
)

type CreateBookInput struct {
	Title     *string
	Author    *string
	Year      *uint16
	Available *bool
}

func CreateBook(ctx context.Context, req CreateBookInput) (model.Book, error) {
	select {
	case <-ctx.Done():
		return model.Book{}, ctx.Err()
	default:
		if req.Author == nil || req.Available == nil || req.Title == nil || req.Year == nil {
			return model.Book{}, apperrors.ErrRequiredFields
		}

		if err := model.ValidateBookFields(req.Title, req.Author, req.Year); err != nil {
			return model.Book{}, err
		}

		book, err := repository.CreateBook(ctx, repository.CreateBookPayload{
			Title:     req.Title,
			Author:    req.Author,
			Year:      req.Year,
			Available: req.Available,
		})

		if err != nil {
			return model.Book{}, err
		}

		return book, nil
	}
}
