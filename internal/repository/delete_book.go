package repository

import (
	"context"
	"library-api/internal/apperrors"
)

func DeleteBookByID(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		mutex.Lock()
		defer mutex.Unlock()
		if ctx.Err() != nil {
			return ctx.Err()
		}
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
}
