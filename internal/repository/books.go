package repository

import (
	"context"
	"library-api/internal/model"
)

type BookRepository struct {
	nextID int
	books  []*model.Book
}

var bookRepository BookRepository = BookRepository{
	nextID: 3,
	books: []*model.Book{
		{
			ID:        1,
			Title:     "1984",
			Author:    "Geroge Orwell",
			Year:      uint16(1956),
			Available: true,
		},
		{
			ID:        1,
			Title:     "Kallocain",
			Author:    "Karin Boye",
			Year:      uint16(1937),
			Available: false,
		},
	},
}

func GetBooks(ctx context.Context) []*model.Book {
	return bookRepository.books
}

type CreateBookRequest struct {
	Title     *string `json:"title"`
	Author    *string `json:"author"`
	Year      *uint16 `json:"year"`
	Available *bool   `json:"available"`
}

func CreateBook(ctx context.Context, req CreateBookRequest) {

	book := model.Book{
		ID:        bookRepository.nextID,
		Title:     *req.Title,
		Author:    *req.Author,
		Year:      *req.Year,
		Available: *req.Available,
	}

	bookRepository.nextID++

	bookRepository.books = append(bookRepository.books, &book)
	// return &book
}
