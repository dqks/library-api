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

func createBookHandler(payload CreateBookPayload) chan model.Book {
	resChan := make(chan model.Book, 1)

	go func() {
		book := model.Book{
			ID:        bookRepository.nextID,
			Title:     *payload.Title,
			Author:    *payload.Author,
			Year:      *payload.Year,
			Available: *payload.Available,
		}

		resChan <- book
	}()

	return resChan
}

func CreateBook(ctx context.Context, payload CreateBookPayload) (model.Book, error) {
	mutex.Lock()
	defer mutex.Unlock()
	select {
	case <-ctx.Done():
		return model.Book{}, ctx.Err()
	case book := <-createBookHandler(payload):
		bookRepository.nextID++
		bookRepository.books = append(bookRepository.books, book)
		return book, nil
	}
}
