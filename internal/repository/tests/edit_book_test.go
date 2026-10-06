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
		name           string
		id             int
		title          string
		author         string
		year           uint16
		available      bool
		titleNil       bool
		authorNil      bool
		yearNil        bool
		availableNil   bool
		wantBook       model.Book
		wantErr        error
		wantCancel     bool
		wantRepoChange bool
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
			wantErr:        nil,
			wantCancel:     false,
			wantRepoChange: true,
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
			wantErr:        nil,
			wantCancel:     false,
			wantRepoChange: true,
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
			wantErr:        nil,
			wantCancel:     false,
			wantRepoChange: true,
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
			wantErr:        nil,
			wantCancel:     false,
			wantRepoChange: true,
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
			wantErr:        nil,
			wantCancel:     false,
			wantRepoChange: true,
		},
		{
			name:           "canceled",
			id:             1,
			title:          "Test title",
			author:         "Test author",
			year:           2000,
			available:      true,
			wantBook:       model.Book{},
			wantErr:        context.Canceled,
			wantCancel:     true,
			wantRepoChange: false,
		},
		{
			name:           "not found",
			id:             3,
			title:          "Test title",
			author:         "Test author",
			year:           2000,
			available:      true,
			titleNil:       false,
			authorNil:      false,
			yearNil:        false,
			availableNil:   false,
			wantBook:       model.Book{},
			wantErr:        apperrors.ErrNotFound,
			wantCancel:     false,
			wantRepoChange: false,
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

			ctxGet := context.WithoutCancel(context.Background())
			booksBefore, err := repo.GetBooks(ctxGet, repository.GetBooksPayload{})

			if err != nil {
				t.Fatal("failed to get books")
			}

			book, err := repo.EditBookByID(ctx, tt.id, payload)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected err %v but got %v", tt.wantErr, err)
			}

			if book != tt.wantBook {
				t.Fatalf("expected book %v but got %v", tt.wantBook, book)
			}

			booksAfter, err := repo.GetBooks(ctxGet, repository.GetBooksPayload{})

			if err != nil {
				t.Fatal("failed to get books")
			}

			if len(booksAfter) != len(booksBefore) {
				t.Fatal("got repo unexpectedly changed")
			}

			for i := range booksAfter {
				if booksAfter[i] != booksBefore[i] && !tt.wantRepoChange {
					t.Fatalf(
						"got book %v changed unexpectedly %v",
						booksBefore[i],
						booksAfter[i],
					)
				} else if tt.wantRepoChange &&
					booksAfter[i] != booksBefore[i] &&
					((tt.wantBook.ID != booksAfter[i].ID) || (tt.wantBook.ID != booksBefore[i].ID)) {
					t.Fatalf(
						"got book %v changed unexpectedly %v",
						booksBefore[i],
						booksAfter[i],
					)
				}
			}

			if !tt.wantCancel && tt.wantErr == nil {
				index := slices.Contains(booksAfter, book)

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
		id         int
	}{
		{
			name:       "success",
			iterations: 100,
			id:         1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := createTestRepo()

			if repo == nil {
				t.Fatal("failed to create repo")
			}

			ctxGet := context.WithoutCancel(context.Background())
			bookBefore, err := repo.GetBookByID(ctxGet, tt.id)

			if err != nil {
				t.Fatalf("failed to get book by id %d", tt.id)
			}

			ctx := context.WithoutCancel(context.Background())
			var wg sync.WaitGroup
			data := make(chan struct {
				year uint16
				book model.Book
				err  error
			}, tt.iterations)

			for i := 1; i <= tt.iterations; i++ {
				wg.Add(1)
				year := uint16(i)
				payload := repository.EditBookPayload{
					Year: &year,
				}

				go func() {
					defer wg.Done()
					book, err := repo.EditBookByID(ctx, tt.id, payload)
					data <- struct {
						year uint16
						book model.Book
						err  error
					}{
						year: *payload.Year,
						book: book,
						err:  err,
					}
				}()
			}

			wg.Wait()

			close(data)

			for d := range data {
				if d.err != nil {
					t.Fatalf("failed to edit book wtih error %v", d.err)
				}
				if d.book.Year != uint16(d.year) {
					t.Fatalf("got book's year %d but expected %d", d.book.Year, d.year)
				}
				if bookBefore.Author != d.book.Author ||
					bookBefore.Title != d.book.Title ||
					bookBefore.Available != d.book.Available ||
					bookBefore.ID != d.book.ID {
					t.Fatalf("book's fields got changed unexpectedly %v", d.book)
				}
			}

			bookAfter, err := repo.GetBookByID(ctxGet, tt.id)

			if err != nil {
				t.Fatalf("failed to get book by id %d", tt.id)
			}

			if bookBefore.Author != bookAfter.Author ||
				bookBefore.Title != bookAfter.Title ||
				bookBefore.Available != bookAfter.Available ||
				bookBefore.ID != bookAfter.ID {
				t.Fatalf("book's fields got changed unexpectedly %v", bookAfter)
			}

			if bookAfter.Year < 1 || bookAfter.Year > 100 {
				t.Fatalf("got book's year out of bounds %d", bookAfter.Year)
			}
		})
	}
}
