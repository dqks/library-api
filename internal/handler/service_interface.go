package handler

import (
	"context"
	"library-api/internal/model"
	"library-api/internal/service"
)

// service.Create возвращает конкретную реализацию
// интерфейс Service принадлежаит потребителю (зависимости) а не реализации

// т.е.
// прописывая интерфейс внутри handler он говорит всем остальным слоям
// "мне похуй что вы мне дадите, мне главное дать структуру с методами ..."
// а service возвращает конкретную реализацию

type Service interface {
	CreateBook(ctx context.Context, req service.CreateBookInput) (model.Book, error)
	DeleteBookByID(ctx context.Context, id int) error
	EditBookByID(ctx context.Context, id int, req service.EditBookInput) (model.Book, error)
	GetBookByID(ctx context.Context, id int) (model.Book, error)
	GetBooks(ctx context.Context, params service.GetBooksQueryParams) ([]model.Book, error)
}
