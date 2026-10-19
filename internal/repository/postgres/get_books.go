package postgres

import (
	"context"
	"fmt"
	"library-api/internal/model"

	repository "library-api/internal/repository"

	"github.com/jackc/pgx/v5"
)

func (r *BookRepository) GetBooks(ctx context.Context, p repository.GetBooksPayload) ([]model.Book, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		books := make([]model.Book, 0)
		var rows pgx.Rows
		var err error
		if p.Available == nil {
			rows, err = r.pool.Query(
				ctx,
				`
				SELECT 
				id, title, author, year, available 
				FROM books
				ORDER BY title
				`,
			)
		} else {
			rows, err = r.pool.Query(
				ctx,
				`
				SELECT 
				id, title, author, year, available 
				FROM books
				WHERE available = $1
				ORDER BY title
				`,
				*p.Available,
			)
		}

		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		for rows.Next() {
			var book model.Book

			err = rows.Scan(
				&book.ID,
				&book.Title,
				&book.Author,
				&book.Year,
				&book.Available,
			)

			if err != nil {
				return nil, fmt.Errorf("%w", err)
			}

			books = append(books, book)
		}

		if rows.Err() != nil {
			return nil, fmt.Errorf("%w", rows.Err())
		}

		return books, nil
	}
}
