package repository_test

import (
	"context"
	"errors"
	"library-api/internal/apperrors"
	"library-api/internal/model"
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
				Title:     "1984",
				Author:    "George Orwell",
				Year:      uint16(1956),
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
			id:       10,
			wantBook: model.Book{},
			wantErr:  context.Canceled,
			cancel:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := createTestRepo()

			if repo == nil {
				t.Fatal("failed to create repository")
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if tt.cancel {
				cancel()
			}

			book, err := repo.GetBookByID(ctx, tt.id)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected err %v but got %v", tt.wantErr, err)
			}

			if book != tt.wantBook {
				t.Fatalf("expected book %v but got %v", tt.wantBook, book)
			}
		})
	}
}
