package repository

import (
	"library-api/internal/model"
	"sync"
)

type BookRepository struct {
	NextID int
	Books  []model.Book
	Mutex  sync.RWMutex
}
