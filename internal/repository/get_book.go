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
		r.Mutex.RLock()
		defer r.Mutex.RUnlock()

		if ctx.Err() != nil {
			return model.Book{}, ctx.Err()
		}

		var index = -1

		for i := range r.Books {
			if r.Books[i].ID == id {
				index = i
				break
			}
		}

		if index == -1 {
			return model.Book{}, apperrors.ErrNotFound
		}

		return r.Books[index], nil
	}
}
