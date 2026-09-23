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
		r.Mutex.Lock()
		defer r.Mutex.Unlock()
		if ctx.Err() != nil {
			return ctx.Err()
		}
		var index = -1
		for i := range r.Books {
			if r.Books[i].ID == id {
				index = i
				break
			}
		}
		if index == -1 {
			return apperrors.ErrNotFound
		}
		r.Books = append(
			r.Books[:index], r.Books[index+1:]...,
		)
		return nil
	}
}
