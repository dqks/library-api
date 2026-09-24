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
	if nextID <= 0 {
		return nil
	}

	for i := range books {
		if books[i].ID >= nextID || books[i].ID <= 0 {
			return nil
		}
	}

	for i := range books {
		for j := range books {
			if books[i].ID == books[j].ID && i != j {
				return nil
			}
		}
	}

	// Мы не должны передавать тот же самый backing array
	// иначе если мы передадим просто books, то вне репозитория
	// мы сможем менять этот слайс из-за backing array
	repoBooks := make([]model.Book, 0, len(books))
	repoBooks = append(repoBooks, books...)

	return &BookRepository{
		nextID: nextID,
		books:  repoBooks,
	}
}
