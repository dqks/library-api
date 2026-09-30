package repository

import (
	"context"
	"errors"
	"library-api/internal/model"
	"slices"
	"sync"
	"testing"
)

func createTestRepo() *BookRepository {
	return CreateRepo(
		3,
		[]model.Book{
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
	)
}

func TestCreateBookConcurrentCancel(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantLen int
		book    model.Book
	}{
		{
			name:    "concurrent test",
			ctx:     context.Background(),
			wantLen: 2,
			book: model.Book{
				ID:        3,
				Title:     "New book",
				Author:    "New author",
				Year:      2025,
				Available: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := createTestRepo()
			if repo == nil {
				t.Fatalf("failed to create test repo")
			}
			testChan := make(chan bool, 1)
			repo.beforeLock = func() {
				testChan <- true
			}
			var err error
			ctx, cancel := context.WithCancel(tt.ctx)
			defer cancel()
			var wg sync.WaitGroup
			wg.Add(3)
			waitChan := make(chan bool, 1)

			go func() {
				defer wg.Done()
				repo.mutex.Lock()
				go func() {
					defer wg.Done()
					_, err = repo.CreateBook(ctx, CreateBookPayload{
						Title:     &tt.book.Title,
						Author:    &tt.book.Author,
						Year:      &tt.book.Year,
						Available: &tt.book.Available,
					})
				}()
				<-testChan
				waitChan <- true
			}()

			go func() {
				defer wg.Done()
				<-waitChan
				cancel()
				repo.mutex.Unlock()
			}()

			wg.Wait()

			if !errors.Is(err, context.Canceled) {
				t.Fatalf("expected %v but got %v", context.Canceled, err)
			}

			ctxGet := context.WithoutCancel(context.Background())
			books, err := repo.GetBooks(ctxGet, GetBooksPayload{})

			if err != nil {
				t.Fatalf("failed to get books")
			}

			if len(books) != tt.wantLen {
				t.Fatalf("got len %d but expected %d", len(books), tt.wantLen)
			}

			if slices.Contains(books, tt.book) {
				t.Fatalf("got book with identical fiends that shouldn't appear %v", tt.book)
			}
		})
	}
}
