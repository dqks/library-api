package postgres

import (
	"context"
	"errors"
	"fmt"
	"library-api/internal/apperrors"
	"library-api/internal/model"

	"github.com/jackc/pgx/v5"
)

func (r *BookRepository) GetUser(ctx context.Context, id int) (model.User, error) {
	select {
	case <-ctx.Done():
		return model.User{}, ctx.Err()
	default:
		var user model.User
		err := r.pool.QueryRow(
			ctx,
			`
			SELECT id, name FROM users WHERE id = $1;
			`,
			id,
		).Scan(
			&user.ID,
			&user.Name,
		)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return model.User{}, fmt.Errorf("%w", apperrors.ErrNotFound)
			}
			return model.User{}, err
		}

		return user, nil
	}
}
