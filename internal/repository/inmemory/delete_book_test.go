package inmemory_test

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/repository"
	"library-api/internal/repository/inmemory"
	"library-api/internal/testutil"
	"sync"
	"testing"
)

func TestDeleteBookByID(t *testing.T) {
	tests := []struct {
		cancel  bool
		name    string
		id      int
		wantErr error
	}{
		{
			name:    "success",
			id:      1,
			wantErr: nil,
			cancel:  false,
		},
		{
			name:    "not found",
			id:      5,
			wantErr: apperrors.ErrNotFound,
			cancel:  false,
		},
		{
			name:    "canceled",
			id:      1,
			wantErr: context.Canceled,
			cancel:  true,
		},
	}

	testBooks := []model.Book{
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
	}

	nextID := 3

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := inmemory.Create(nextID, testBooks)

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

			if tt.cancel {
				cancel()
			}

			fields := testutil.DeleteBookFields{
				ID:          tt.id,
				Ctx:         ctx,
				Cancel:      tt.cancel,
				BooksBefore: booksBefore,
				Repo:        repo,
				WantErr:     tt.wantErr,
			}

			testutil.CheckDeleteBook(t, &fields)
		})
	}
}

func TestDeleteBookByIDMultiple(t *testing.T) {
	tests := []struct {
		name            string
		ids             []int
		wantDeleteCount int8
		wantErrs        map[error]int
	}{
		{
			name:            "success",
			ids:             []int{1, 4, 6},
			wantDeleteCount: 3,
			wantErrs: map[error]int{
				nil: 3,
			},
		},
		{
			name:            "error, not found",
			ids:             []int{1, 1, 2},
			wantDeleteCount: 2,
			wantErrs: map[error]int{
				apperrors.ErrNotFound: 1,
				nil:                   2,
			},
		},
	}

	testBooks := []model.Book{
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
			Available: true,
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
	}
	nextID := 7

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := inmemory.Create(
				nextID,
				testBooks,
			)

			if repo == nil {
				t.Fatal("failed to create repository")
			}
			errChan := make(chan error, len(tt.ids))
			var wg sync.WaitGroup
			ctx := context.WithoutCancel(context.Background())
			booksBefore, err := repo.GetBooks(ctx, repository.GetBooksPayload{})

			if err != nil {
				t.Fatal("failed to get books")
			}

			for _, id := range tt.ids {
				wg.Add(1)
				go func() {
					defer wg.Done()
					err := repo.DeleteBookByID(ctx, id)
					errChan <- err
				}()
			}

			wg.Wait()

			close(errChan)

			gotErrs := make(map[error]int)

			for e := range errChan {
				gotErrs[e] = gotErrs[e] + 1
			}

			for key, val := range tt.wantErrs {
				if errCount, ok := gotErrs[key]; ok {
					if val != errCount {
						t.Fatalf("expected errCount %d of error %v but got %d", val, key, errCount)
					}
				} else {
					t.Fatalf("expected to have err %v but didn't", key)
				}
			}

			booksAfter, err := repo.GetBooks(ctx, repository.GetBooksPayload{})

			if err != nil {
				t.Fatal("failed to get books")
			}

			if len(booksBefore)-int(tt.wantDeleteCount) != len(booksAfter) {
				t.Fatal("failed to delete book")
			}

			for i := range tt.ids {
				for j := range booksAfter {
					if tt.ids[i] == booksAfter[j].ID {
						t.Fatalf("failed to delete book by id %d", tt.ids[i])
					}
				}
			}
		})
	}
}
