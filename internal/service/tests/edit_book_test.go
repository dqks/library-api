package tests

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/service"
	"testing"
	"time"
)

func TestEditBookByID(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		author    string
		year      uint16
		available bool
		cancel    bool
		id        int
		wantBook  model.Book
		wantErr   error
	}{
		{
			name:      "success",
			title:     "New Title",
			author:    "New Author",
			year:      uint16(2000),
			available: true,
			cancel:    false,
			wantBook: model.Book{
				ID:        1,
				Title:     "New Title",
				Author:    "New Author",
				Year:      uint16(2000),
				Available: true,
			},
			wantErr: nil,
		},
		{
			name:      "err, empty string title",
			title:     "",
			author:    "New Author",
			year:      uint16(0),
			available: true,
			cancel:    false,
			wantBook:  model.Book{},
			wantErr:   apperrors.ErrInvalidValues,
		},
		{
			name:      "err, empty string author",
			title:     "",
			author:    "New Author",
			year:      uint16(0),
			available: true,
			cancel:    false,
			wantBook:  model.Book{},
			wantErr:   apperrors.ErrInvalidValues,
		},
		{
			name:      "err, year more than now",
			title:     "New Title",
			author:    "New Author",
			year:      uint16(time.Now().Year() + 1),
			available: true,
			cancel:    false,
			wantBook:  model.Book{},
			wantErr:   apperrors.ErrInvalidValues,
		},
		{
			name:      "err, canceled",
			title:     "New Title",
			author:    "New Author",
			year:      uint16(time.Now().Year()),
			available: true,
			cancel:    true,
			wantBook:  model.Book{},
			wantErr:   context.Canceled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if tt.cancel {
				cancel()
			}

			s := service.Create(&MockRepo{})

			book, err := s.EditBookByID(ctx, tt.id, service.EditBookInput{
				Title:     &tt.title,
				Author:    &tt.author,
				Year:      &tt.year,
				Available: &tt.available,
			})

			if book != tt.wantBook {
				t.Fatalf("expected book %v but got %v", tt.wantBook, book)
			}

			if err != tt.wantErr {
				t.Fatalf("expected err %v but got %v", tt.wantErr, err)
			}
		})
	}
}
