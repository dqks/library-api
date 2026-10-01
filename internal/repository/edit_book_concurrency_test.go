package repository

import (
	"context"
	"errors"
	"library-api/internal/model"
	"sync"
	"testing"
)

func TestEditBookConcurrenyCancel(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		author    string
		year      uint16
		available bool
		id        int
		book      model.Book
	}{
		{
			name:      "cancel",
			title:     "Test book",
			author:    "Test author",
			year:      2000,
			available: true,
			id:        1,
			book: model.Book{
				ID:        1,
				Title:     "1984",
				Author:    "George Orwell",
				Year:      uint16(1956),
				Available: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := createTestRepo()

			if repo == nil {
				t.Fatalf("failed to create repo")
			}

			var wg sync.WaitGroup
			var err error
			wg.Add(3)
			testChan := make(chan bool)
			waitChan := make(chan bool)

			repo.beforeLock = func() {
				waitChan <- true
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func() {
				defer wg.Done()

				<-testChan
				go func() {
					defer wg.Done()
					_, err = repo.EditBookByID(ctx, tt.id, EditBookPayload{
						Title:     &tt.title,
						Author:    &tt.author,
						Year:      &tt.year,
						Available: &tt.available,
					})
				}()

				<-waitChan

				testChan <- true

			}()

			go func() {
				defer wg.Done()
				repo.mutex.Lock()
				testChan <- true
				<-testChan
				cancel()
				repo.mutex.Unlock()
			}()

			wg.Wait()

			if !errors.Is(err, context.Canceled) {
				t.Fatalf("got err %v but expected %v", err, context.Canceled)
			}

			ctxGet := context.WithoutCancel(context.Background())
			book, err := repo.GetBookByID(ctxGet, tt.id)

			if err != nil {
				t.Fatalf("failed to get books")
			}

			if book != tt.book {
				t.Fatalf("got book %v but expected %v", book, tt.book)
			}
		})
	}
}
