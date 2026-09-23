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
		r.Mutex.Lock()
		defer r.Mutex.Unlock()

		if ctx.Err() != nil {
			return model.Book{}, ctx.Err()
		}

		var index = -1

		for i := range r.Books {
			if r.Books[i].ID == id {
				index = i
				break
			}
		}

		if index == -1 {
			return model.Book{}, apperrors.ErrNotFound
		}

		if p.Author != nil {
			r.Books[index].Author = *p.Author
		}

		if p.Available != nil {
			r.Books[index].Available = *p.Available
		}

		if p.Title != nil {
			r.Books[index].Title = *p.Title
		}

		if p.Year != nil {
			r.Books[index].Year = *p.Year
		}

		return r.Books[index], nil

	}
}
