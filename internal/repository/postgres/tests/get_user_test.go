package postgres_tests

import (
	"context"
	"errors"
	"library-api/internal/apperrors"
	"library-api/internal/database"
	"library-api/internal/model"
	"testing"
	"time"
)

func TestGetUser(t *testing.T) {
	tests := []struct {
		name            string
		cancel          bool
		deadlineContext bool
		id              int
		wantErr         error
		wantUser        model.User
	}{
		{
			name: "success",
			id:   1,
			wantUser: model.User{
				ID:   1,
				Name: "user 1",
			},
		},
		{
			name:    "canceled",
			id:      1,
			cancel:  true,
			wantErr: context.Canceled,
		},
		{
			name:            "deadline exceeded",
			id:              1,
			deadlineContext: true,
			wantErr:         context.DeadlineExceeded,
		},
		{
			name:    "not found",
			id:      1000,
			wantErr: apperrors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, pool := CreatePostgresTestRepo(t)

			err := database.UpMigrations(pool)
			if err != nil {
				t.Fatalf("failed to apply up migrations %v", err)
			}

			t.Cleanup(func() {
				err := database.DownMigrations(pool)
				if err != nil {
					t.Fatalf("failed to apply down migrations %v", err)
				}
			})

			SeedUsers(t, pool)

			var ctx context.Context
			var cancel context.CancelFunc

			if tt.deadlineContext {
				ctx, cancel = context.WithDeadline(context.Background(), time.Time{})
				defer cancel()

			} else {
				ctx, cancel = context.WithCancel(context.Background())
				defer cancel()
			}

			if tt.cancel {
				cancel()
			}

			user, err := repo.GetUser(ctx, tt.id)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"expected err to be %v but got %v",
					tt.wantErr,
					err,
				)
			}

			if user != tt.wantUser {
				t.Fatalf(
					"expected user to be %v but got %v",
					tt.wantUser,
					user,
				)
			}

		})
	}
}
