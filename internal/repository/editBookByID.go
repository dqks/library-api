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

func EditBookByID(ctx context.Context, id int, p EditBookPayload) (model.Book, error) {
	defer mutex.Unlock()
	select {
	case <-ctx.Done():
		return model.Book{}, ctx.Err()
	default:
		mutex.Lock()
		var index = -1

		for i := range bookRepository.books {
			if bookRepository.books[i].ID == id {
				index = i
				break
			}
		}

		if index == -1 {
			return model.Book{}, apperrors.ErrNotFound
		}

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
	}
}
