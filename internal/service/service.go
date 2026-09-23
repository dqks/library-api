package service

import "library-api/internal/repository"

type BookService struct {
	Repo *repository.BookRepository
}
