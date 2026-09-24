package service

import "library-api/internal/repository"

type BookService struct {
	repo *repository.BookRepository
}

func CreateBookService(repo *repository.BookRepository) *BookService {
	return &BookService{repo: repo}
}
