package postgres_test

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/repository"
	"library-api/internal/testutil"
	"testing"
)

func TestDeleteBook(t *testing.T) {
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
			id:      6,
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

			ctxGet := context.WithoutCancel(context.Background())
			booksBefore, err := repo.GetBooks(ctxGet, repository.GetBooksPayload{})

			if err != nil {
				t.Fatal("failed to get books")
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
