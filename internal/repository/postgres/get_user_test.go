package postgres_test

import (
	"context"
	"errors"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/testutil"
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
			repo, pool := testutil.CreatePostgresTestRepo(t)
			t.Cleanup(func() {
				pool.Close()
			})
			testutil.SeedUsers(t, pool)

			t.Cleanup(func() {
				ctx := context.Background()
				_, err := testutil.ClearUsers(ctx, pool)
				if err != nil {
					t.Errorf("failed to clear users %v", err)
				}
			})

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
