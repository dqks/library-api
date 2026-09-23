package repository

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
)

func GetBookByID(ctx context.Context, id int) (model.Book, error) {
	select {
	case <-ctx.Done():
		return model.Book{}, ctx.Err()
	default:
		mutex.RLock()
		defer mutex.RUnlock()

		if ctx.Err() != nil {
			return model.Book{}, ctx.Err()
		}

		var index = -1

		for i := range bookRepository.books {
			if bookRepository.books[i].ID == id {
				index = i
				break
			}
		}

		if index == -1 {
			return model.Book{}, apperrors.ErrNotFound
		}

		return bookRepository.books[index], nil
	}
}
