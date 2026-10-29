package postgres_tests

import (
	"context"
	"errors"
	"library-api/internal/apperrors"
	"library-api/internal/database"
	"library-api/internal/model"
	"library-api/internal/repository"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {

	tests := []struct {
		name            string
		userName        *string
		cancel          bool
		deadlineContext bool
		wantErr         error
		wantUser        model.User
	}{
		{
			name:     "success",
			userName: new("test name"),
			wantUser: model.User{
				ID:   6,
				Name: "test name",
			},
		},
		{
			name:     "success, len 3",
			userName: new("tet"),
			wantUser: model.User{
				ID:   6,
				Name: "tet",
			},
		},
		{
			name:     "success, len 25",
			userName: new("fsdgfgsdfgsdfgsdhghdаывав"),
			wantUser: model.User{
				ID:   6,
				Name: "fsdgfgsdfgsdfgsdhghdаывав",
			},
		},
		{
			name:     "err, len 2",
			userName: new("te"),
			wantErr:  apperrors.ErrUserCreate,
		},
		{
			name:     "err, len 26",
			userName: new("fsdgfgsdfgsdfgsdhghdаывавl"),
			wantErr:  apperrors.ErrUserCreate,
		},
		{
			name:     "err, empty string",
			userName: new(""),
			wantErr:  apperrors.ErrUserCreate,
		},
		{
			name:     "canceled",
			userName: new("test name"),
			cancel:   true,
			wantErr:  context.Canceled,
		},
		{
			name:            "deadline exceeded",
			deadlineContext: true,
			userName:        new("test name"),
			wantErr:         context.DeadlineExceeded,
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
				t.Fatalf("failed to apply migrations")
			}
			t.Cleanup(func() {
				err := database.DownMigrations(pool)
				if err != nil {
					t.Errorf("failed to apply migrations")
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

			user, err := repo.CreateUser(ctx, repository.CreateUserPayload{
				Name: tt.userName,
			})

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"expected err to be %v but got %v",
					tt.wantErr,
					err,
				)
			}

			if user != tt.wantUser {
				t.Fatalf(
					"expected user to be %v\nbut got %v",
					tt.wantUser,
					user,
				)
			}
		})
	}
}
