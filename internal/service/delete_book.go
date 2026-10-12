package service

import (
	"context"
)

func (s *bookService) DeleteBookByID(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return s.repo.DeleteBookByID(ctx, id)
	}
}
