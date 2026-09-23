package repository

import (
	"context"
	"library-api/internal/apperrors"
)

func (r *BookRepository) DeleteBookByID(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		r.mutex.Lock()
		defer r.mutex.Unlock()
		if ctx.Err() != nil {
			return ctx.Err()
		}
		var index = -1
		for i := range r.books {
			if r.books[i].ID == id {
				index = i
				break
			}
		}
		if index == -1 {
			return apperrors.ErrNotFound
		}
		r.books = append(
			r.books[:index], r.books[index+1:]...,
		)
		return nil
	}
}
