package service

import "library-api/internal/repository"

type BookService struct {
	repo *repository.BookRepository
}

func CreateBookService(repo *repository.BookRepository) *BookService {

	if repo == nil {
		return nil
	}

	return &BookService{repo: repo}
}
