package repository

import (
	"context"
	"library-api/internal/model"
)

type GetBooksPayload struct {
	Available *bool
}

func getBooksHandler(params GetBooksPayload) chan []model.Book {
	resChan := make(chan []model.Book)

	go func() {
		mutex.RLock()
		defer mutex.RUnlock()
		if params.Available != nil {
			availableBooks := make([]model.Book, 0, len(bookRepository.books))
			for i := range bookRepository.books {
				if bookRepository.books[i].Available == *params.Available {
					availableBooks = append(availableBooks, bookRepository.books[i])
				}
			}
			resChan <- availableBooks
			return
		}

		books := make([]model.Book, 0, len(bookRepository.books))
		books = append(books, bookRepository.books...)
	}()

	return resChan
}

func GetBooks(ctx context.Context, params GetBooksPayload) ([]model.Book, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case books := <-getBooksHandler(params):
		return books, nil
	}
}
