package service

import (
	"context"
)

func (s *BookService) DeleteBookByID(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return s.Repo().DeleteBookByID(ctx, id)
	}
}
