package postgres

import (
	"context"
	"errors"
	"fmt"
	"library-api/internal/apperrors"
	"library-api/internal/model"

	"github.com/jackc/pgx/v5"
)

func (r *BookRepository) GetBookByID(ctx context.Context, id int) (model.Book, error) {
	select {
	case <-ctx.Done():
		return model.Book{}, ctx.Err()
	default:
		var book model.Book

		err := r.pool.QueryRow(
			ctx,
			`
			SELECT id, title, author, year, available
			FROM books
			WHERE id = $1
			`,
			id,
		).Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Year,
			&book.Available,
		)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return model.Book{}, fmt.Errorf("%w", apperrors.ErrNotFound)
			}

			return model.Book{}, fmt.Errorf("%w", err)

		}

		return book, nil
	}
}
