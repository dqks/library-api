package postgres_tests

import (
	"context"
	"library-api/internal/database"
	"library-api/internal/model"
	"library-api/internal/testutil"
	"testing"
)

func TestCreateBook(t *testing.T) {
	tests := []struct {
		name         string
		cancel       bool
		ctx          context.Context
		title        string
		author       string
		year         uint16
		available    bool
		wantError    error
		wantBook     model.Book
		wantToAppear bool
		wantLen      int
	}{
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
			cancel:       true,
			wantLen:      5,
		},
		{
			name:      "success",
			ctx:       context.Background(),
			title:     "New book",
			author:    "New author",
			year:      2025,
			available: true,
			wantError: nil,
			wantBook: model.Book{
				ID:        6,
				Title:     "New book",
				Author:    "New author",
				Year:      2025,
				Available: true,
			},
			wantToAppear: true,
			cancel:       false,
			wantLen:      6,
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

			fields := testutil.CreateFields{
				Title:        tt.title,
				Author:       tt.author,
				Year:         tt.year,
				Available:    tt.available,
				Ctx:          ctx,
				Repo:         repo,
				WantError:    tt.wantError,
				WantBook:     tt.wantBook,
				WantToAppear: tt.wantToAppear,
				WantLen:      tt.wantLen,
			}

			testutil.CheckCreateBook(t, &fields)
		})
	}
}
