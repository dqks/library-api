package tests

import (
	"context"
	"errors"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/service"
	"library-api/internal/testutil"
	"testing"
)

func TestGetUser(t *testing.T) {
	tests := []struct {
		name            string
		cancel          bool
		deadlineContext bool
		id              int
		//repo settings
		repoUser model.User
		repoErr  error
		//assert
		wantUser model.User
		wantErr  error
	}{
		{
			name: "success",
			id:   1,
			repoUser: model.User{
				ID:   1,
				Name: "test",
			},
			repoErr: nil,
			wantUser: model.User{
				ID:   1,
				Name: "test",
			},
			wantErr: nil,
		},
		{
			name:    "not found",
			id:      1,
			repoErr: apperrors.ErrNotFound,
			wantErr: apperrors.ErrNotFound,
		},
		{
			name:    "canceled",
			cancel:  true,
			id:      1,
			wantErr: context.Canceled,
		},
		{
			name:            "deadline exceeded",
			deadlineContext: true,
			id:              1,
			wantErr:         context.DeadlineExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := MockRepo{
				getUserErr: tt.repoErr,
				user:       tt.repoUser,
			}

			s := service.CreateUserService(&repo)

			ctx, cancel := testutil.GetNeededCtx(tt.deadlineContext)
			t.Cleanup(func() {
				cancel()
			})

			if tt.cancel {
				cancel()
			}

			user, err := s.GetUser(ctx, tt.id)

			if tt.wantUser != user {
				t.Fatalf(
					"expected user to be %v but got %v",
					tt.wantUser,
					user,
				)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"expected user to be %v but got %v",
					tt.wantUser,
					user,
				)
			}

			if repo.gotInside && repo.getUserID != tt.id {
				t.Fatalf(
					"expected to call repo with id %d but got %d",
					tt.id,
					repo.getUserID,
				)
			}
		})
	}
}
