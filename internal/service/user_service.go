package service

import (
	"context"
	"library-api/internal/model"
	"library-api/internal/repository"
)

type UserRepository interface {
	CreateUser(ctx context.Context, payload repository.CreateUserPayload) (model.User, error)
	GetUser(ctx context.Context, id int) (model.User, error)
}

type userService struct {
	repo UserRepository
}

func CreateUserService(repo UserRepository) *userService {
	if repo == nil {
		return nil
	}

	return &userService{repo: repo}
}
