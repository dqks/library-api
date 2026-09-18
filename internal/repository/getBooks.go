package repository

import (
	"context"
	"library-api/internal/model"
)

type GetBooksPayload struct {
	Available *bool
}

func GetBooks(ctx context.Context, params GetBooksPayload) ([]model.Book, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		mutex.RLock()
		defer mutex.RUnlock()
		if params.Available != nil {
			availableBooks := make([]model.Book, 0, len(bookRepository.books))
			for i := range bookRepository.books {
				if bookRepository.books[i].Available == *params.Available {
					availableBooks = append(availableBooks, bookRepository.books[i])
				}
			}
			return availableBooks, nil
		}

		// Чтобы не возвращать тот же самый backing array
		// и подавить concurrency проблему
		books := make([]model.Book, 0, len(bookRepository.books))
		books = append(books, bookRepository.books...)
		return books, nil
	}
}
