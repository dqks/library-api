package service

import (
	"context"
	"library-api/internal/model"
)

func (s *BookService) GetBookByID(ctx context.Context, id int) (model.Book, error) {
	select {
	case <-ctx.Done():
		return model.Book{}, ctx.Err()
	default:
		return s.Repo().GetBookByID(ctx, id)
	}
}
