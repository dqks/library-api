package repository

import (
	"library-api/internal/model"
	"sync"
)

type BookRepository struct {
	nextID int
	books  []model.Book
}

var bookRepository = BookRepository{
	nextID: 3,
	books: []model.Book{
		{
			ID:        1,
			Title:     "1984",
			Author:    "George Orwell",
			Year:      uint16(1956),
			Available: true,
		},
		{
			ID:        2,
			Title:     "Kallocain",
			Author:    "Karin Boye",
			Year:      uint16(1937),
			Available: false,
		},
	},
}

var mutex sync.RWMutex
