package tests

import (
	"context"
	"errors"
	"library-api/internal/apperrors"
	"library-api/internal/service"
	"testing"
)

func TestLoanBook(t *testing.T) {
	tests := []struct {
		//test set ups
		name   string
		bookID int
		userID int
		cancel bool
		//repo set up
		userErr     error
		loanBookErr error
		// result
		wantErr error
	}{
		{
			name:    "success",
			bookID:  1,
			userID:  1,
			wantErr: nil,
		},
		{
			name:    "canceled",
			cancel:  true,
			wantErr: context.Canceled,
		},
		{
			name:    "err, user not found",
			bookID:  1,
			userID:  1,
			userErr: apperrors.ErrNotFound,
			wantErr: apperrors.ErrNotFound,
		},
		{
			name:        "err, loan error",
			bookID:      1,
			userID:      1,
			loanBookErr: apperrors.ErrInternal,
			wantErr:     apperrors.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockRepo{
				getUserErr:  tt.userErr,
				loanBookErr: tt.loanBookErr,
			}
			service := service.CreateLoanService(repo)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if tt.cancel {
				cancel()
			}

			err := service.LoanBook(ctx, tt.bookID, tt.userID)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf(
					"expected loan err to be %v but got %v",
					tt.wantErr,
					err,
				)
			}

			if repo.getUserID != tt.userID {
				t.Errorf(
					"expected getUserID to be %d but got %d",
					tt.userID,
					repo.userID,
				)
			}

			if repo.getUserErr == nil && repo.err == nil {
				if repo.loanBookID != tt.bookID {
					t.Errorf(
						"expected loanBookID to be %d but got %d",
						tt.bookID,
						repo.loanBookID,
					)
				}

				if repo.loanUserID != tt.userID && repo.getUserErr == nil && repo.err == nil {
					t.Errorf(
						"expected loanUserID to be %d but got %d",
						tt.userID,
						repo.loanUserID,
					)
				}
			}
		})
	}
}
