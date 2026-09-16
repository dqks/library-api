package service

import (
	"context"
	"library-api/internal/model"
	"library-api/internal/repository"
)

func GetBooks(ctx context.Context, params repository.GetBooksQueryParams) []*model.Book {
	select {
	case books := <-func() chan []*model.Book {
		resultChan := make(chan []*model.Book)
		go func() {
			resultChan <- repository.GetBooks(ctx, params)
		}()
		return resultChan
	}():
		return books
	case <-ctx.Done():
		return nil
	}
}
