package repository

import (
	"context"
	"library-api/internal/model"
)

type GetBooksPayload struct {
	Available *bool
}

func (r *BookRepository) GetBooks(ctx context.Context, params GetBooksPayload) ([]model.Book, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		r.mutex.RLock()
		defer r.mutex.RUnlock()

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		if params.Available != nil {
			availableBooks := make([]model.Book, 0, len(r.books))
			for i := range r.books {
				if r.books[i].Available == *params.Available {
					availableBooks = append(availableBooks, r.books[i])
				}
			}
			return availableBooks, nil
		}

		books := make([]model.Book, 0, len(r.books))
		books = append(books, r.books...)
		return books, nil
	}
}
