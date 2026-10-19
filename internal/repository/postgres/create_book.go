package postgres

import (
	"context"
	"fmt"
	"library-api/internal/model"
	repository "library-api/internal/repository"
)

func (r *BookRepository) CreateBook(ctx context.Context, p repository.CreateBookPayload) (model.Book, error) {
	select {
	case <-ctx.Done():
		return model.Book{}, ctx.Err()
	default:
		var book model.Book
		err := r.pool.QueryRow(ctx,
			`
			INSERT INTO books (title, author, available, year)
			VALUES ($1, $2, $3, $4)
			RETURNING id, title, author, available, year
			`,
			p.Title,
			p.Author,
			p.Available,
			p.Year,
		).Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Available,
			&book.Year,
		)
		if err != nil {
			return model.Book{}, fmt.Errorf(
				"ошибка создания книги: %w",
				err,
			)
		}
		return book, nil
	}
}
