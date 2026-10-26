package service

type userService struct {
	repo UserRepository
}

func CreateUserService(repo UserRepository) *userService {
	if repo == nil {
		return nil
	}

	return &userService{repo: repo}
}
