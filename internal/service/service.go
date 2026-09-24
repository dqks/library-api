package service

import "library-api/internal/repository"

type BookService struct {
	repo *repository.BookRepository
}

func CreateBookService(repo *repository.BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (b *BookService) Repo() *repository.BookRepository {
	if b.repo == nil {
		return &repository.BookRepository{}
	}

	return b.repo
}
