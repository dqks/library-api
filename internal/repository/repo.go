package repository

import (
	"library-api/internal/model"
	"sync"
)

type BookRepository struct {
	nextID int
	books  []model.Book
	mutex  sync.RWMutex
}

func CreateRepo(nextID int, books []model.Book) *BookRepository {
	return &BookRepository{
		nextID: nextID,
		books:  books,
	}
}
