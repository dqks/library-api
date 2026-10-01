package repository_test

import (
	"context"
	"errors"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/repository"
	"sync"
	"testing"
)

func TestDeleteBookByID(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		wantErr    error
		wantCancel bool
	}{
		{
			name:       "success",
			id:         1,
			wantErr:    nil,
			wantCancel: false,
		},
		{
			name:       "not found",
			id:         5,
			wantErr:    apperrors.ErrNotFound,
			wantCancel: false,
		},
		{
			name:       "canceled",
			id:         1,
			wantErr:    context.Canceled,
			wantCancel: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := createTestRepo()

			if repo == nil {
				t.Fatal("failed to create repository")
			}

			ctxGet := context.WithoutCancel(context.Background())
			booksBefore, err := repo.GetBooks(ctxGet, repository.GetBooksPayload{})

			if err != nil {
				t.Fatal("failed to get books")
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if tt.wantCancel {
				cancel()
			}

			err = repo.DeleteBookByID(ctx, tt.id)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected err %v but got %v", tt.wantErr, err)
			}

			books, err := repo.GetBooks(ctxGet, repository.GetBooksPayload{})
			if tt.wantErr != nil {
				for i := 0; i < len(books); i++ {
					if books[i] != booksBefore[i] {
						t.Fatalf("got changes in books after delete")
					}
				}

				if len(books) != len(booksBefore) {
					t.Fatalf(
						"got len of books %d but before delete it was %d",
						len(books),
						len(booksBefore),
					)
				}
			}

			if tt.wantErr == nil && !tt.wantCancel {

				if err != nil {
					t.Fatal("failed to get books")
				}

				for _, b := range books {
					if b.ID == tt.id {
						t.Fatalf("expected to delete book %v but got it in books", b)
					}
				}
			}
		})
	}
}

func TestDeleteBookByIDMultiple(t *testing.T) {
	tests := []struct {
		name    string
		ids     []int
		wantErr error
	}{
		{
			name:    "success",
			ids:     []int{1, 4, 6},
			wantErr: nil,
		},
		{
			name:    "same ids",
			ids:     []int{1, 1},
			wantErr: apperrors.ErrNotFound,
		},
	}

	createTestRepo := func() *repository.BookRepository {
		return repository.CreateRepo(
			7,
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
				{
					ID:        3,
					Title:     "Kallocain",
					Author:    "Karin Boye",
					Year:      uint16(1937),
					Available: false,
				},
				{
					ID:        4,
					Title:     "Kallocain",
					Author:    "Karin Boye",
					Year:      uint16(1937),
					Available: false,
				},
				{
					ID:        5,
					Title:     "Kallocain",
					Author:    "Karin Boye",
					Year:      uint16(1937),
					Available: false,
				},
				{
					ID:        6,
					Title:     "Kallocain",
					Author:    "Karin Boye",
					Year:      uint16(1937),
					Available: false,
				},
			},
		)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := createTestRepo()

			if repo == nil {
				t.Fatal("failed to create repository")
			}
			errChan := make(chan error, len(tt.ids))
			var wg sync.WaitGroup
			ctx := context.WithoutCancel(context.Background())

			for _, id := range tt.ids {
				wg.Add(1)
				go func() {
					defer wg.Done()
					err := repo.DeleteBookByID(ctx, id)
					if err != nil {
						errChan <- err
					}
				}()
			}

			wg.Wait()

			close(errChan)
			for e := range errChan {
				if !errors.Is(e, tt.wantErr) {
					t.Errorf("expected err %v but got %v", tt.wantErr, e)
				}
			}

			books, err := repo.GetBooks(ctx, repository.GetBooksPayload{})

			if err != nil {
				t.Fatal("failed to get books")
			}

			for i := range tt.ids {
				for j := range books {
					if tt.ids[i] == books[j].ID {
						t.Fatalf("failed to delete book by id %d", tt.ids[i])
					}
				}
			}
		})
	}
}
