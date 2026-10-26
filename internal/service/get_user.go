package service

import (
	"context"
	"library-api/internal/model"
)

func (s *userService) GetUser(ctx context.Context, id int) (model.User, error) {
	select {
	case <-ctx.Done():
		return model.User{}, ctx.Err()
	default:
		return s.repo.GetUser(ctx, id)
	}
}
