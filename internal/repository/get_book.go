package repository

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
)

func getBookByIDHandler(id int) chan model.Book {
	resChan := make(chan model.Book, 1)

	go func() {
		mutex.RLock()
		defer mutex.RUnlock()
		var index = -1

		for i := range bookRepository.books {
			if bookRepository.books[i].ID == id {
				index = i
				break
			}
		}

		if index != -1 {
			resChan <- bookRepository.books[index]
		} else {
			resChan <- model.Book{}
		}

	}()

	return resChan
}

func GetBookByID(ctx context.Context, id int) (model.Book, error) {
	select {
	case <-ctx.Done():
		return model.Book{}, ctx.Err()
	case book := <-getBookByIDHandler(id):
		if book.ID != 0 {
			return book, nil
		} else {
			return book, apperrors.ErrNotFound
		}
	}
}
