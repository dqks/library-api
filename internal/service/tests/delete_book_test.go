package tests

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/service"
	"testing"
)

func TestDeleteBookByID(t *testing.T) {
	tests := []struct {
		name           string
		id             int
		cancel         bool
		repositoryBook model.Book
		repositoryErr  error
		wantErr        error
	}{
		{
			name:    "success",
			id:      1,
			cancel:  false,
			wantErr: nil,
		},
		{
			name:    "canceled by service",
			id:      1,
			cancel:  true,
			wantErr: context.Canceled,
		},
		{
			name:          "canceled by repository",
			id:            1,
			wantErr:       context.Canceled,
			repositoryErr: context.Canceled,
		},
		{
			name:          "not found",
			id:            1,
			wantErr:       apperrors.ErrNotFound,
			repositoryErr: apperrors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if tt.cancel {
				cancel()
			}

			s := service.Create(&MockRepo{
				err: tt.repositoryErr,
			})

			err := s.DeleteBookByID(ctx, tt.id)

			if err != tt.wantErr {
				t.Fatalf("expected err %v but got %v", tt.wantErr, err)
			}
		})
	}
}
