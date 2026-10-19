package pg

import (
	"context"
	"fmt"
	"library-api/internal/apperrors"
)

func (r *BookRepository) DeleteBookByID(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		tag, err := r.pool.Exec(
			ctx,
			"DELETE FROM books WHERE id = $1",
			id,
		)

		if err != nil {
			return fmt.Errorf("ошибка удаления книги: %w", err)
		}

		if tag.RowsAffected() == 0 {
			return apperrors.ErrNotFound
		}

		return nil
	}
}
