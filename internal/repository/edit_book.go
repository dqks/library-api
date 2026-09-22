package repository

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
)

type EditBookPayload struct {
	Title     *string
	Author    *string
	Year      *uint16
	Available *bool
}

func editBookHandler(id int, p EditBookPayload) chan model.Book {
	resChan := make(chan model.Book, 1)

	go func() {
		mutex.Lock()
		defer mutex.Unlock()
		var index = -1

		for i := range bookRepository.books {
			if bookRepository.books[i].ID == id {
				index = i
				break
			}
		}

		if index != -1 {
			if p.Author != nil {
				bookRepository.books[index].Author = *p.Author
			}

			if p.Available != nil {
				bookRepository.books[index].Available = *p.Available
			}

			if p.Title != nil {
				bookRepository.books[index].Title = *p.Title
			}

			if p.Year != nil {
				bookRepository.books[index].Year = *p.Year
			}

			resChan <- bookRepository.books[index]
		} else {
			resChan <- model.Book{}
		}
	}()

	return resChan
}

func EditBookByID(ctx context.Context, id int, p EditBookPayload) (model.Book, error) {
	select {
	case <-ctx.Done():
		return model.Book{}, ctx.Err()
	case book := <-editBookHandler(id, p):
		if book.ID != 0 {
			return book, nil
		} else {
			return book, apperrors.ErrNotFound
		}
	}
}
