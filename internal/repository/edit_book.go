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

func editBookHandler(id int) chan int {
	resChan := make(chan int)

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

		resChan <- index
	}()

	return resChan
}

func EditBookByID(ctx context.Context, id int, p EditBookPayload) (model.Book, error) {
	select {
	case <-ctx.Done():
		return model.Book{}, ctx.Err()
	case index := <-editBookHandler(id):
		if index != -1 {
			mutex.Lock()
			defer mutex.Unlock()
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

			return bookRepository.books[index], nil
		} else {
			return model.Book{}, apperrors.ErrNotFound
		}
	}
}
