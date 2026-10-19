package service

import (
	"context"
	"library-api/internal/model"
	repository "library-api/internal/repository/inmemory"
)

type GetBooksQueryParams struct {
	Available *bool
}

func (s *bookService) GetBooks(ctx context.Context, params GetBooksQueryParams) ([]model.Book, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.GetBooks(ctx, repository.GetBooksPayload{Available: params.Available})
	}
}
