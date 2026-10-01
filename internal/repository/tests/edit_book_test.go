package repository_test

import (
	"context"
	"errors"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/repository"
	"slices"
	"sync"
	"testing"
)

func TestEditBookByID(t *testing.T) {
	tests := []struct {
		name         string
		id           int
		title        string
		author       string
		year         uint16
		available    bool
		titleNil     bool
		authorNil    bool
		yearNil      bool
		availableNil bool
		wantBook     model.Book
		wantErr      error
		wantCancel   bool
	}{
		{
			name:         "success all fields",
			id:           1,
			title:        "Test title",
			author:       "Test author",
			year:         2000,
			available:    true,
			titleNil:     false,
			authorNil:    false,
			yearNil:      false,
			availableNil: false,
			wantBook: model.Book{
				ID:        1,
				Title:     "Test title",
				Author:    "Test author",
				Year:      2000,
				Available: true,
			},
			wantErr:    nil,
			wantCancel: false,
		},
		{
			name:         "success only title",
			id:           1,
			title:        "Test title",
			author:       "",
			year:         0,
			available:    false,
			titleNil:     false,
			authorNil:    true,
			yearNil:      true,
			availableNil: true,
			wantBook: model.Book{
				ID:        1,
				Title:     "Test title",
				Author:    "George Orwell",
				Year:      uint16(1956),
				Available: true,
			},
			wantErr:    nil,
			wantCancel: false,
		},
		{
			name:         "success only author",
			id:           1,
			title:        "",
			author:       "Test author",
			year:         0,
			available:    false,
			titleNil:     true,
			authorNil:    false,
			yearNil:      true,
			availableNil: true,
			wantBook: model.Book{
				ID:        1,
				Title:     "1984",
				Author:    "Test author",
				Year:      uint16(1956),
				Available: true,
			},
			wantErr:    nil,
			wantCancel: false,
		},
		{
			name:         "success only year",
			id:           1,
			title:        "",
			author:       "",
			year:         2000,
			available:    false,
			titleNil:     true,
			authorNil:    true,
			yearNil:      false,
			availableNil: true,
			wantBook: model.Book{
				ID:        1,
				Title:     "1984",
				Author:    "George Orwell",
				Year:      uint16(2000),
				Available: true,
			},
			wantErr:    nil,
			wantCancel: false,
		},
		{
			name:         "success only available",
			id:           1,
			title:        "",
			author:       "",
			year:         2000,
			available:    false,
			titleNil:     true,
			authorNil:    true,
			yearNil:      true,
			availableNil: false,
			wantBook: model.Book{
				ID:        1,
				Title:     "1984",
				Author:    "George Orwell",
				Year:      uint16(1956),
				Available: false,
			},
			wantErr:    nil,
			wantCancel: false,
		},
		{
			name:       "canceled",
			id:         1,
			title:      "Test title",
			author:     "Test author",
			year:       2000,
			available:  true,
			wantBook:   model.Book{},
			wantErr:    context.Canceled,
			wantCancel: true,
		},
		{
			name:         "not found",
			id:           3,
			title:        "Test title",
			author:       "Test author",
			year:         2000,
			available:    true,
			titleNil:     false,
			authorNil:    false,
			yearNil:      false,
			availableNil: false,
			wantBook:     model.Book{},
			wantErr:      apperrors.ErrNotFound,
			wantCancel:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := createTestRepo()

			if repo == nil {
				t.Fatal("failed to create repository")
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if tt.wantCancel {
				cancel()
			}

			payload := repository.EditBookPayload{}

			if !tt.titleNil {
				payload.Title = &tt.title
			}

			if !tt.authorNil {
				payload.Author = &tt.author
			}

			if !tt.yearNil {
				payload.Year = &tt.year
			}

			if !tt.availableNil {
				payload.Available = &tt.available
			}

			book, err := repo.EditBookByID(ctx, tt.id, payload)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected err %v but got %v", tt.wantErr, err)
			}

			if book != tt.wantBook {
				t.Fatalf("expected book %v but got %v", tt.wantBook, book)
			}

			ctxGet := context.WithoutCancel(context.Background())
			books, err := repo.GetBooks(ctxGet, repository.GetBooksPayload{})

			if err != nil {
				t.Fatal("failed to get books")
			}

			if !tt.wantCancel && tt.wantErr == nil {
				index := slices.Contains(books, book)

				if !index {
					t.Fatalf("failed to find book in repository %v", book)
				}
			}
		})
	}
}

func TestEditBookByIDConcurrent(t *testing.T) {
	tests := []struct {
		name       string
		iterations int
	}{
		{
			name:       "success",
			iterations: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := createTestRepo()

			if repo == nil {
				t.Fatal("failed to create repo")
			}

			ctx := context.WithoutCancel(context.Background())
			yearSlice := make([]uint16, 0, 100)
			var wg sync.WaitGroup
			var errEdit error

			for i := 1; i <= tt.iterations; i++ {
				wg.Add(1)
				year := uint16(i)
				payload := repository.EditBookPayload{
					Year: &year,
				}

				go func() {
					defer wg.Done()
					_, err := repo.EditBookByID(ctx, 1, payload)

					if err != nil {
						errEdit = err
					} else {
						yearSlice = append(yearSlice, *payload.Year)
					}
				}()
			}

			wg.Wait()

			if errEdit != nil {
				t.Fatal("failed to edit book")
			}

			ctxGet := context.WithoutCancel(context.Background())
			books, err := repo.GetBooks(ctxGet, repository.GetBooksPayload{})

			if err != nil {
				t.Fatalf("failed to get books")
			}

			if len(yearSlice) != tt.iterations {
				t.Fatalf("got %d edits but expected %d", len(yearSlice), tt.iterations)
			}

			for i := range yearSlice {
				for j := range yearSlice {
					if yearSlice[i] == yearSlice[j] && i != j {
						t.Fatalf("got not unique year in edits")
					}
				}
			}

			if books[0].Year != yearSlice[len(yearSlice)-1] {
				t.Fatalf(
					"expected final year %d but got %d",
					books[0].Year,
					yearSlice[len(yearSlice)-1],
				)
			}
		})
	}
}
