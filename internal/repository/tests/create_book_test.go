package repository_test

import (
	"context"
	"errors"
	"library-api/internal/model"
	"library-api/internal/repository"
	"sync"
	"testing"
)

func createTestRepo() *repository.BookRepository {
	return repository.CreateRepo(
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

func TestCreateBook(t *testing.T) {
	tests := []struct {
		name         string
		ctx          context.Context
		title        string
		author       string
		year         uint16
		available    bool
		wantError    error
		wantBook     model.Book
		wantToAppear bool
		wantToCancel bool
		wantLen      int
	}{
		{
			name:      "success",
			ctx:       context.Background(),
			title:     "New book",
			author:    "New author",
			year:      2025,
			available: true,
			wantError: nil,
			wantBook: model.Book{
				ID:        3,
				Title:     "New book",
				Author:    "New author",
				Year:      2025,
				Available: true,
			},
			wantToAppear: true,
			wantToCancel: false,
			wantLen:      3,
		},
		{
			name:         "canceled",
			ctx:          context.Background(),
			title:        "New book",
			author:       "New author",
			year:         2025,
			available:    true,
			wantError:    context.Canceled,
			wantBook:     model.Book{},
			wantToAppear: false,
			wantToCancel: true,
			wantLen:      2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := repository.CreateBookPayload{
				Title:     &tt.title,
				Author:    &tt.author,
				Year:      &tt.year,
				Available: &tt.available,
			}

			repo := createTestRepo()

			if repo == nil {
				t.Fatalf("failed to create test repo")
			}

			ctx, cancel := context.WithCancel(tt.ctx)
			defer cancel()

			if tt.wantToCancel {
				cancel()
			}

			book, err := repo.CreateBook(ctx, payload)

			if !errors.Is(err, tt.wantError) {
				t.Fatalf("got error %v but expected %v", err, tt.wantError)
			}

			if book != tt.wantBook {
				t.Fatalf("got %v but expected %v", book, tt.wantBook)
			}

			ctxGet := context.WithoutCancel(context.Background())

			books, err := repo.GetBooks(ctxGet, repository.GetBooksPayload{})

			if err != nil {
				t.Fatalf("failed to get books")
			}

			if (books[len(books)-1] != tt.wantBook && tt.wantToAppear) ||
				(books[len(books)-1] == tt.wantBook && !tt.wantToAppear) {
				t.Fatalf("new book %v didn't appear on book list", book)
			}

			if len(books) != tt.wantLen {
				t.Fatalf("got len %d but expercted %d", len(books), tt.wantLen)
			}
		})
	}
}

func TestCreateBookConcurrent(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		title     string
		author    string
		year      uint16
		available bool
		wantID    int
		wantLen   int
	}{
		{
			name:      "concurrent test",
			ctx:       context.Background(),
			title:     "New book",
			author:    "New author",
			year:      2025,
			available: true,
			wantID:    102,
			wantLen:   102,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := createTestRepo()

			if repo == nil {
				t.Fatalf("failed to create test repo")
			}

			ctx := context.WithoutCancel(tt.ctx)
			var wg sync.WaitGroup

			for i := 0; i < 100; i++ {
				wg.Add(1)

				go func() {
					defer wg.Done()

					repo.CreateBook(ctx, repository.CreateBookPayload{
						Title:     &tt.title,
						Author:    &tt.author,
						Year:      &tt.year,
						Available: &tt.available,
					})
				}()
			}

			wg.Wait()

			ctxGet := context.WithoutCancel(context.Background())
			books, err := repo.GetBooks(ctxGet, repository.GetBooksPayload{})

			if err != nil {
				t.Fatalf("failed to get books")
			}

			if books[len(books)-1].ID != tt.wantID {
				t.Fatalf("got id %d but expected %d", books[len(books)-1].ID, tt.wantID)
			}

			if len(books) != tt.wantLen {
				t.Fatalf("got len %d but expected %d", len(books), tt.wantLen)
			}

			for i := range books {
				for j := range books {
					if books[i].ID == books[j].ID && i != j {
						t.Fatalf("got not unique id %d", books[i].ID)
					}
				}
			}
		})
	}
}
