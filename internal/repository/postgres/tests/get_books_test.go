package postgres_tests

import (
	"context"
	"library-api/internal/database"
	"library-api/internal/model"
	"library-api/internal/repository"
	"library-api/internal/testutil"
	"testing"
)

func TestGetBooks(t *testing.T) {
	tests := []struct {
		cancel       bool
		name         string
		available    bool
		availableNil bool
		wantErr      error
		wantBooks    []model.Book
	}{
		{
			name:         "success, no payload",
			available:    false,
			availableNil: true,
			wantErr:      nil,
			wantBooks: []model.Book{
				{
					ID:        1,
					Title:     "Title 1",
					Author:    "Author 1",
					Year:      uint16(2000),
					Available: true,
				},
				{
					ID:        2,
					Title:     "Title 2",
					Author:    "Author 2",
					Year:      uint16(2000),
					Available: true,
				},
				{
					ID:        3,
					Title:     "Title 3",
					Author:    "Author 3",
					Year:      uint16(2000),
					Available: false,
				},
				{
					ID:        4,
					Title:     "Title 4",
					Author:    "Author 4",
					Year:      uint16(2000),
					Available: true,
				},
				{
					ID:        5,
					Title:     "Title 5",
					Author:    "Author 5",
					Year:      uint16(2000),
					Available: true,
				},
			},
		},
		{
			name:         "success, payload with available true",
			available:    true,
			availableNil: false,
			wantErr:      nil,
			wantBooks: []model.Book{
				{
					ID:        1,
					Title:     "Title 1",
					Author:    "Author 1",
					Year:      uint16(2000),
					Available: true,
				},
				{
					ID:        2,
					Title:     "Title 2",
					Author:    "Author 2",
					Year:      uint16(2000),
					Available: true,
				},
				{
					ID:        4,
					Title:     "Title 4",
					Author:    "Author 4",
					Year:      uint16(2000),
					Available: true,
				},
				{
					ID:        5,
					Title:     "Title 5",
					Author:    "Author 5",
					Year:      uint16(2000),
					Available: true,
				},
			},
		},
		{
			name:         "success, payload with available false",
			available:    false,
			availableNil: false,
			wantErr:      nil,
			wantBooks: []model.Book{
				{
					ID:        3,
					Title:     "Title 3",
					Author:    "Author 3",
					Year:      uint16(2000),
					Available: false,
				},
			},
		},
		{
			name:         "canceled",
			available:    false,
			availableNil: true,
			wantErr:      context.Canceled,
			cancel:       true,
			wantBooks:    []model.Book{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, pool := CreatePostgresTestRepo(t)
			t.Cleanup(func() {
				pool.Close()
			})

			err := database.UpMigrations(pool)
			if err != nil {
				t.Fatalf("ошибка применения миграций: %v", err)
			}

			t.Cleanup(func() {
				err := database.DownMigrations(pool)
				if err != nil {
					t.Errorf("ошибка применения миграций: %v", err)
				}
			})

			SeedBooks(t, pool)

			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(func() {
				cancel()
			})

			if tt.cancel {
				cancel()
			}

			payload := repository.GetBooksPayload{}

			if !tt.availableNil {
				payload.Available = &tt.available
			}

			books, err := repo.GetBooks(ctx, payload)

			fields := testutil.GetBooksFields{
				Err:       err,
				Books:     books,
				Cancel:    tt.cancel,
				Payload:   payload,
				Ctx:       ctx,
				Repo:      repo,
				WantErr:   tt.wantErr,
				WantBooks: tt.wantBooks,
			}

			testutil.CheckGetBooks(t, &fields)
		})
	}
}
