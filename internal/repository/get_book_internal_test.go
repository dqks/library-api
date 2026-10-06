package repository

import (
	"context"
	"library-api/internal/model"
	"sync"
	"testing"
)

func TestGetBookByIDReaderWriter(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		author    string
		year      uint16
		available bool
		id        int
		wantBook  model.Book
	}{
		{
			name:      "get new version of book",
			title:     "Test title",
			author:    "Test author",
			year:      2000,
			available: false,
			id:        1,
			wantBook: model.Book{
				ID:        1,
				Title:     "Test title",
				Author:    "Test author",
				Year:      2000,
				Available: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := createTestRepo()

			if repo == nil {
				t.Fatal("failed to create repository")
			}

			testChan := make(chan bool)
			repo.afterLock = func() {
				testChan <- true
			}

			var book model.Book
			ctx := context.WithoutCancel(context.Background())
			var wg sync.WaitGroup
			wg.Add(3)

			go func() {
				defer wg.Done()
				<-testChan
				book, _ = repo.GetBookByID(ctx, tt.id)
			}()

			go func() {
				defer wg.Done()
				go func() {
					defer wg.Done()
					repo.EditBookByID(ctx, tt.id, EditBookPayload{
						Title:     &tt.title,
						Author:    &tt.author,
						Year:      &tt.year,
						Available: &tt.available,
					})
				}()
			}()

			wg.Wait()

			if book != tt.wantBook {
				t.Fatalf("expected book %v but got %v", tt.wantBook, book)
			}
		})
	}
}
