package postgres_tests

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/repository"
	"library-api/internal/testutil"
	"testing"
)

func TestEditBookByID(t *testing.T) {
	tests := []struct {
		name           string
		cancel         bool
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
			cancel:         false,
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
			cancel:         false,
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
			cancel:         false,
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
			cancel:         false,
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
			cancel:         false,
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
			cancel:         true,
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
			cancel:         false,
			wantRepoChange: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, pool := CreatePostgresTestRepo(t)

			SeedBooks(t, pool)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if tt.cancel {
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

			fields := testutil.EditFields{
				Err:            err,
				Book:           book,
				CtxGet:         ctxGet,
				Repo:           repo,
				Cancel:         tt.cancel,
				BooksBefore:    booksBefore,
				WantRepoChange: tt.wantRepoChange,
				WantErr:        tt.wantErr,
				WantBook:       tt.wantBook,
			}

			testutil.CheckEdit(t, &fields)
		})
	}
}
