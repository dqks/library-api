package repository

import (
	"context"
	"fmt"
	"library-api/internal/model"
	"sync"
	"testing"
	"time"
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
		name      string
		ctx       context.Context
		title     string
		author    string
		year      uint16
		available bool
		stopIndex int
	}{
		{
			name:      "concurrent test",
			ctx:       context.Background(),
			title:     "New book",
			author:    "New author",
			year:      2025,
			available: true,
			stopIndex: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := createTestRepo()

			if repo == nil {
				t.Fatalf("failed to create test repo")
			}

			ctx, cancel := context.WithCancel(tt.ctx)
			defer cancel()
			var wg sync.WaitGroup
			wg.Add(2)
			waitChan := make(chan bool, 1)

			go func() {
				defer wg.Done()
				repo.mutex.Lock()
				waitChan <- true
				cancel()
				repo.mutex.Unlock()
			}()

			go func() {
				defer wg.Done()
				<-waitChan
				repo.CreateBook(ctx, CreateBookPayload{
					Title:     &tt.title,
					Author:    &tt.author,
					Year:      &tt.year,
					Available: &tt.available,
				})
			}()

			wg.Wait()

			time.Sleep(5 * time.Second)

			ctxGet := context.WithoutCancel(context.Background())
			books, err := repo.GetBooks(ctxGet, GetBooksPayload{})

			if err != nil {
				t.Fatalf("failed to get books")
			}

			fmt.Println(books)
		})
	}
}
