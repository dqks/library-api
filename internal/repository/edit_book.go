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

func (r *BookRepository) EditBookByID(ctx context.Context, id int, p EditBookPayload) (model.Book, error) {
	select {
	case <-ctx.Done():
		return model.Book{}, ctx.Err()
	default:
		r.mutex.Lock()
		defer r.mutex.Unlock()

		if ctx.Err() != nil {
			return model.Book{}, ctx.Err()
		}

		var index = -1

		for i := range r.books {
			if r.books[i].ID == id {
				index = i
				break
			}
		}

		if index == -1 {
			return model.Book{}, apperrors.ErrNotFound
		}

		if p.Author != nil {
			r.books[index].Author = *p.Author
		}

		if p.Available != nil {
			r.books[index].Available = *p.Available
		}

		if p.Title != nil {
			r.books[index].Title = *p.Title
		}

		if p.Year != nil {
			r.books[index].Year = *p.Year
		}

		return r.books[index], nil

	}
}
