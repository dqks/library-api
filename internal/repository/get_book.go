package repository

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
)

func (r *BookRepository) GetBookByID(ctx context.Context, id int) (model.Book, error) {
	select {
	case <-ctx.Done():
		return model.Book{}, ctx.Err()
	default:
		r.mutex.RLock()
		defer r.mutex.RUnlock()

		if ctx.Err() != nil {
			return model.Book{}, ctx.Err()
		}

		var index = -1

		for i := range r.books {
			if r.books[i].ID == id {
				index = i
				break
			}
		}

		if index == -1 {
			return model.Book{}, apperrors.ErrNotFound
		}

		return r.books[index], nil
	}
}
