package repository

import (
	"context"
	"library-api/internal/apperrors"
)

func deleteBookHandler(id int) chan int {
	resChan := make(chan int)

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

		resChan <- index
	}()

	return resChan
}

func DeleteBookByID(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case index := <-deleteBookHandler(id):
		if index != -1 {
			mutex.Lock()
			defer mutex.Unlock()
			bookRepository.books = append(
				bookRepository.books[:index], bookRepository.books[index+1:]...,
			)
			return nil
		} else {
			return apperrors.ErrNotFound
		}
	}
}
