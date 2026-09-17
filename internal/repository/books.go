package repository

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"sync"
)

type BookRepository struct {
	nextID int
	books  []*model.Book
}

var bookRepository = BookRepository{
	nextID: 3,
	books: []*model.Book{
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

var mutex sync.Mutex

type GetBooksPayload struct {
	Available bool
}

func GetBooks(ctx context.Context, params GetBooksPayload) []*model.Book {
	mutex.Lock()
	defer mutex.Unlock()

	if params.Available {
		availableBooks := make([]*model.Book, 0, len(bookRepository.books))
		index := 0
		for i := range bookRepository.books {
			if bookRepository.books[i].Available {
				availableBooks = append(availableBooks, bookRepository.books[i])
				index++
			}
		}
		return availableBooks
	}

	return bookRepository.books
}

type CreateBookPayload struct {
	Title     *string `json:"title"`
	Author    *string `json:"author"`
	Year      *uint16 `json:"year"`
	Available *bool   `json:"available"`
}

func CreateBook(ctx context.Context, payload CreateBookPayload) *model.Book {
	mutex.Lock()
	defer mutex.Unlock()

	book := model.Book{
		ID:        bookRepository.nextID,
		Title:     *payload.Title,
		Author:    *payload.Author,
		Year:      *payload.Year,
		Available: *payload.Available,
	}

	bookRepository.nextID++

	bookRepository.books = append(bookRepository.books, &book)

	return &book
}

func DeleteBookByID(ctx context.Context, id int) error {
	mutex.Lock()
	defer mutex.Unlock()
	var index = -1

	for i := range bookRepository.books {
		if bookRepository.books[i].ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return apperrors.ErrNotFound
	}

	bookRepository.books = append(
		bookRepository.books[:index], bookRepository.books[index+1:]...,
	)

	return nil
}

func GetBookByID(ctx context.Context, id int) (*model.Book, error) {
	mutex.Lock()
	defer mutex.Unlock()
	var index = -1

	for i := range bookRepository.books {
		if bookRepository.books[i].ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return nil, apperrors.ErrNotFound
	}

	return bookRepository.books[index], nil
}

type EditBookPayload struct {
	Title     *string `json:"title"`
	Author    *string `json:"author"`
	Year      *uint16 `json:"year"`
	Available *bool   `json:"available"`
}

func EditBookByID(ctx context.Context, id int, p EditBookPayload) (*model.Book, error) {
	mutex.Lock()
	defer mutex.Unlock()

	var index = -1

	for i := range bookRepository.books {
		if bookRepository.books[i].ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return nil, apperrors.ErrNotFound
	}

	if p.Author != nil {
		bookRepository.books[index].Author = *p.Author
	}

	if p.Available != nil {
		bookRepository.books[index].Available = *p.Available
	}

	if p.Title != nil {
		bookRepository.books[index].Title = *p.Title
	}

	if p.Year != nil {
		bookRepository.books[index].Year = *p.Year
	}

	return bookRepository.books[index], nil
}
