package pgrepo

import (
	"context"
	"fmt"
	"library-api/internal/model"
	repository "library-api/internal/repository"
)

func (r *BookRepository) EditBookByID(ctx context.Context, id int, p repository.EditBookPayload) (model.Book, error) {
	select {
	case <-ctx.Done():
		return model.Book{}, ctx.Err()
	default:
		var book model.Book

		err := r.pool.QueryRow(
			ctx,
			`
			UPDATE books
			SET title = COALESCE($2, title),
				author = COALESCE($3, author),
				year = COALESCE($4, year),
				available = COALESCE($5, available)
			WHERE id = $1
			RETURNING id, title, author, year, available;
			`,
			id,
			p.Title,
			p.Author,
			p.Year,
			p.Available,
		).Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Year,
			&book.Available,
		)

		if err != nil {
			return model.Book{}, fmt.Errorf("%w", err)
		}

		return book, nil
	}
}
