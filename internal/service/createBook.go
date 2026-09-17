package service

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/repository"
)

type CreateBookRequest struct {
	Title     *string `json:"title"`
	Author    *string `json:"author"`
	Year      *uint16 `json:"year"`
	Available *bool   `json:"available"`
}

func CreateBook(ctx context.Context, req CreateBookRequest) (model.Book, error) {
	select {
	case <-ctx.Done():
		return model.Book{}, apperrors.ErrContext
	default:
		if err := model.ValidateBookFields(req.Title, req.Author, req.Year, req.Available); err != nil {
			return model.Book{}, err
		}

		book := repository.CreateBook(ctx, repository.CreateBookPayload{
			Title:     req.Title,
			Author:    req.Author,
			Year:      req.Year,
			Available: req.Available,
		})
		return book, nil
	}
}
