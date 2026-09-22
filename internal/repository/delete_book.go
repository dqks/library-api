package repository

import (
	"context"
	"library-api/internal/apperrors"
)

func deleteBookHandler(id int) chan error {
	resChan := make(chan error, 1)

	go func() {
		mutex.Lock()
		defer mutex.Unlock()
		var index = -1

		for i := range bookRepository.books {
			if bookRepository.books[i].ID == id {
				index = i
				break
			}
		}

		if index != -1 {
			bookRepository.books = append(
				bookRepository.books[:index], bookRepository.books[index+1:]...,
			)
			resChan <- nil
		} else {
			resChan <- apperrors.ErrNotFound
		}

	}()

	return resChan
}

func DeleteBookByID(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-deleteBookHandler(id):
		return err
	}
}
