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
		r.Mutex.RLock()
		defer r.Mutex.RUnlock()

		if params.Available != nil {
			availableBooks := make([]model.Book, 0, len(r.Books))
			for i := range r.Books {
				if r.Books[i].Available == *params.Available {
					availableBooks = append(availableBooks, r.Books[i])
				}
			}
			return availableBooks, nil
		}

		books := make([]model.Book, 0, len(r.Books))
		books = append(books, r.Books...)
		return books, nil
	}
}
