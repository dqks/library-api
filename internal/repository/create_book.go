package repository

import (
	"context"
	"fmt"
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
		r.mutex.Lock()
		defer r.mutex.Unlock()

		fmt.Println("unlock inside create")

		if ctx.Err() != nil {
			return model.Book{}, ctx.Err()
		}

		book := model.Book{
			ID:        r.nextID,
			Title:     *payload.Title,
			Author:    *payload.Author,
			Year:      *payload.Year,
			Available: *payload.Available,
		}

		r.nextID++
		r.books = append(r.books, book)
		return book, nil
	}
}
