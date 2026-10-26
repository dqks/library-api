package service

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/repository"
)

type CreateUserInput struct {
	Name *string
}

func (s *userService) CreateUser(ctx context.Context, input CreateUserInput) (model.User, error) {
	select {
	case <-ctx.Done():
		return model.User{}, ctx.Err()
	default:
		if err := model.ValidateUserFields(input.Name); err != nil {
			return model.User{}, apperrors.ErrInvalidValues
		}

		return s.repo.CreateUser(ctx, repository.CreateUserPayload{Name: input.Name})
	}
}
