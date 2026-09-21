package repository

import (
	"context"
	"library-api/internal/model"
)

func getBookByIDHandler(id int) chan int {
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

func GetBookByID(ctx context.Context, id int) (model.Book, error) {
	select {
	case <-ctx.Done():
		return model.Book{}, ctx.Err()
	case index := <-getBookByIDHandler(id):
		mutex.RLock()
		defer mutex.RUnlock()
		return bookRepository.books[index], nil
	}
}
