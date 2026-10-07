package tests

import (
	"context"
	"library-api/internal/model"
	"library-api/internal/repository"
)

type MockRepo struct{}

func (r *MockRepo) CreateBook(ctx context.Context, p repository.CreateBookPayload) (model.Book, error) {
	return model.Book{
		ID:        1,
		Title:     *p.Title,
		Author:    *p.Author,
		Year:      *p.Year,
		Available: *p.Available,
	}, nil
}

func (r *MockRepo) DeleteBookByID(ctx context.Context, id int) error {
	return nil
}

// TODO в тесте и здесь продумать nil поля
// нужен ли на это тест
// но я думаю что нужен, т.к. сервис может взаимодействовать
// с nil полями
func (r *MockRepo) EditBookByID(ctx context.Context, id int, p repository.EditBookPayload) (model.Book, error) {
	return model.Book{
		ID:        1,
		Title:     *p.Title,
		Author:    *p.Author,
		Year:      *p.Year,
		Available: *p.Available,
	}, nil
}

func (r *MockRepo) GetBookByID(ctx context.Context, id int) (model.Book, error) {
	return model.Book{
		ID:        1,
		Title:     "Mock Title",
		Author:    "Mock Author",
		Year:      2000,
		Available: false,
	}, nil
}

func (r *MockRepo) GetBooks(ctx context.Context, params repository.GetBooksPayload) ([]model.Book, error) {
	return []model.Book{
		{
			ID:        1,
			Title:     "Mock Title",
			Author:    "Mock Author",
			Year:      2000,
			Available: false,
		},
	}, nil
}
