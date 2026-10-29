package postgres

import (
	"context"
	"fmt"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	repository "library-api/internal/repository"
)

func (r *BookRepository) CreateUser(ctx context.Context, payload repository.CreateUserPayload) (model.User, error) {
	select {
	case <-ctx.Done():
		return model.User{}, ctx.Err()
	default:
		var user model.User
		err := r.pool.QueryRow(
			ctx,
			`
			INSERT INTO users (name)
			VALUES ($1)
			RETURNING id, name
			`,
			payload.Name,
		).Scan(
			&user.ID,
			&user.Name,
		)

		if err != nil {
			return model.User{}, fmt.Errorf(
				"%w: %v",
				apperrors.ErrUserCreate,
				err,
			)
		}

		return user, nil
	}
}
