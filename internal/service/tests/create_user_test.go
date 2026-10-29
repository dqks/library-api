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

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name            string
		cancel          bool
		deadlineContext bool
		repoUser        model.User
		repoErr         error
		userName        *string
		wantErr         error
		wantUser        model.User
	}{
		{
			name:     "success",
			userName: new("test name"),
			repoUser: model.User{
				ID:   1,
				Name: "test name",
			},
			repoErr: nil,
			wantUser: model.User{
				ID:   1,
				Name: "test name",
			},
		},
		{
			name:     "success, name 3 chars",
			userName: new("test name"),
			repoUser: model.User{
				ID:   1,
				Name: "tes",
			},
			repoErr: nil,
			wantUser: model.User{
				ID:   1,
				Name: "tes",
			},
		},
		{
			name:     "success, name 25 chars",
			userName: new("test name"),
			repoUser: model.User{
				ID:   1,
				Name: "dasdasdasdasdasddasdasdsd",
			},
			repoErr: nil,
			wantUser: model.User{
				ID:   1,
				Name: "dasdasdasdasdasddasdasdsd",
			},
		},
		{
			name:     "canceled",
			cancel:   true,
			userName: new("test name"),
			wantErr:  context.Canceled,
		},
		{
			name:            "deadline",
			deadlineContext: true,
			userName:        new("test name"),
			wantErr:         context.DeadlineExceeded,
		},
		{
			name:     "err, empty name",
			userName: new(""),
			wantErr:  apperrors.ErrInvalidValues,
		},
		{
			name:     "err, name 2 chars",
			userName: new("dd"),
			wantErr:  apperrors.ErrInvalidValues,
		},
		{
			name:     "err, name 26 chars",
			userName: new("dasdasdasdasdasddasdasddsd"),
			wantErr:  apperrors.ErrInvalidValues,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repo := MockRepo{
				err:  tt.repoErr,
				user: tt.repoUser,
			}

			s := service.CreateUserService(&repo)

			ctx, cancel := testutil.GetNeededCtx(tt.deadlineContext)
			t.Cleanup(func() {
				cancel()
			})

			if tt.cancel {
				cancel()
			}

			user, err := s.CreateUser(ctx, service.CreateUserInput{
				Name: tt.userName,
			})

			if user != tt.wantUser {
				t.Fatalf(
					"expected user to be %v but got %v",
					user,
					tt.wantUser,
				)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"expected err to be %v but got %v",
					err,
					tt.wantErr,
				)
			}

			if repo.createUserPayload.Name != tt.userName && repo.gotInside {
				t.Errorf(
					"expected repo's name payload to be %v but got %v",
					tt.userName,
					repo.createUserPayload.Name,
				)
			}

		})
	}
}
