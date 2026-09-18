package repository

import (
	"context"
	"library-api/internal/apperrors"
)

func DeleteBookByID(ctx context.Context, id int) error {
	defer mutex.Unlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		mutex.Lock()
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
