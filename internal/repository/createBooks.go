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

func CreateBook(ctx context.Context, payload CreateBookPayload) (model.Book, error) {
	defer mutex.Unlock()
	select {
	case <-ctx.Done():
		return model.Book{}, ctx.Err()
	default:
		mutex.Lock()
		book := model.Book{
			ID:        bookRepository.nextID,
			Title:     *payload.Title,
			Author:    *payload.Author,
			Year:      *payload.Year,
			Available: *payload.Available,
		}

		bookRepository.nextID++

		bookRepository.books = append(bookRepository.books, book)

		return book, nil
	}
}
