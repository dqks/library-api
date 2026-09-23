package repository

import (
	"context"
	"library-api/internal/model"
)

type CreateBookPayload struct {
	Title     *string
	Author    *string
	Year      *uint16
	Available *bool
}

func (r *BookRepository) CreateBook(ctx context.Context, payload CreateBookPayload) (model.Book, error) {
	select {
	case <-ctx.Done():
		return model.Book{}, ctx.Err()
	default:
		r.Mutex.Lock()
		defer r.Mutex.Unlock()

		if ctx.Err() != nil {
			return model.Book{}, ctx.Err()
		}

		book := model.Book{
			ID:        r.NextID,
			Title:     *payload.Title,
			Author:    *payload.Author,
			Year:      *payload.Year,
			Available: *payload.Available,
		}

		r.NextID++
		r.Books = append(r.Books, book)
		return book, nil
	}
}
