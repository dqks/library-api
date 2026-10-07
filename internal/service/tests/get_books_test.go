package tests

import (
	"context"
	"library-api/internal/model"
	"library-api/internal/service"
	"testing"
)

func TestGetBooks(t *testing.T) {
	tests := []struct {
		name         string
		available    bool
		availableNil bool
		cancel       bool
		wantBooks    []model.Book
		wantErr      error
	}{
		{
			name:      "success, available false",
			available: false,
			cancel:    false,
			wantBooks: []model.Book{
				{
					ID:        1,
					Title:     "Mock Title",
					Author:    "Mock Author",
					Year:      2000,
					Available: false,
				},
			},
			wantErr: nil,
		},
		{
			name:         "success, available true",
			available:    true,
			availableNil: true,
			cancel:       false,
			wantBooks: []model.Book{
				{
					ID:        1,
					Title:     "Mock Title",
					Author:    "Mock Author",
					Year:      2000,
					Available: false,
				},
			},
			wantErr: nil,
		},
		{
			name:      "canceled",
			available: false,
			wantBooks: []model.Book{},
			wantErr:   context.Canceled,
			cancel:    true,
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

			params := service.GetBooksQueryParams{}

			if !tt.availableNil {
				params.Available = &tt.available
			}

			books, err := s.GetBooks(ctx, params)

			if err != tt.wantErr {
				t.Fatalf(
					"expected err %v but got %v",
					tt.wantErr,
					err,
				)
			}

			if len(books) != len(tt.wantBooks) {
				t.Fatalf(
					"expected len of books %d but got %d",
					len(tt.wantBooks),
					len(books),
				)
			}

			for i := range books {
				if books[i] != tt.wantBooks[i] {
					t.Fatalf(
						"expected book %v but got %v",
						tt.wantBooks[i],
						books[i],
					)
				}
			}
		})
	}
}
