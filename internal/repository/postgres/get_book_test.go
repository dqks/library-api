package postgres_test

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/testutil"
	"testing"
)

func TestGetBookByID(t *testing.T) {
	tests := []struct {
		cancel   bool
		name     string
		id       int
		wantBook model.Book
		wantErr  error
	}{
		{
			name: "success",
			id:   1,
			wantBook: model.Book{
				ID:        1,
				Title:     "Title 1",
				Author:    "Author 1",
				Year:      uint16(2000),
				Available: true,
			},
			wantErr: nil,
			cancel:  false,
		},
		{
			name:     "not found",
			id:       10,
			wantBook: model.Book{},
			wantErr:  apperrors.ErrNotFound,
			cancel:   false,
		},
		{
			name:     "canceled",
			id:       1,
			wantBook: model.Book{},
			wantErr:  context.Canceled,
			cancel:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, pool := testutil.CreatePostgresTestRepo(t)
			t.Cleanup(func() {
				pool.Close()
			})

			testutil.SeedBooks(t, pool)

			t.Cleanup(func() {
				ctx := context.Background()
				_, err := testutil.ClearBooks(ctx, pool)
				if err != nil {
					t.Errorf("failed to clear books %v", err)
				}
			})

			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(func() {
				cancel()
			})

			if tt.cancel {
				cancel()
			}

			fields := testutil.GetBookFields{
				ID:       tt.id,
				Ctx:      ctx,
				Repo:     repo,
				WantErr:  tt.wantErr,
				WantBook: tt.wantBook,
			}

			testutil.CheckGetBook(t, &fields)
		})
	}
}
