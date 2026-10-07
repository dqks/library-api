package tests

import (
	"context"
	"library-api/internal/model"
	"library-api/internal/service"
	"testing"
)

func TestGetBook(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		cancel   bool
		wantBook model.Book
		wantErr  error
	}{
		{
			name:   "success",
			id:     1,
			cancel: false,
			wantBook: model.Book{
				ID:        1,
				Title:     "Mock Title",
				Author:    "Mock Author",
				Year:      2000,
				Available: false,
			},
			wantErr: nil,
		},
		{
			name:     "canceled",
			id:       1,
			wantBook: model.Book{},
			wantErr:  context.Canceled,
			cancel:   true,
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

			book, err := s.GetBookByID(ctx, tt.id)

			if err != tt.wantErr {
				t.Fatalf(
					"expected err %v but got %v",
					tt.wantErr,
					err,
				)
			}

			if tt.wantBook != book {
				t.Fatalf(
					"expected book %v but got %v",
					tt.wantBook,
					book,
				)
			}
		})
	}
}
