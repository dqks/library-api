package repository_test

import (
	"context"
	"errors"
	"library-api/internal/model"
	"library-api/internal/repository"
	"testing"
)

func TestGetBooks(t *testing.T) {
	tests := []struct {
		name         string
		available    bool
		availableNil bool
		wantErr      error
		wantCancel   bool
		wantBooks    []model.Book
	}{
		{
			name:         "success, no payload",
			available:    false,
			availableNil: true,
			wantErr:      nil,
			wantBooks: []model.Book{
				{
					ID:        1,
					Title:     "1984",
					Author:    "George Orwell",
					Year:      uint16(1956),
					Available: true,
				},
				{
					ID:        2,
					Title:     "Kallocain",
					Author:    "Karin Boye",
					Year:      uint16(1937),
					Available: false,
				},
				{
					ID:        3,
					Title:     "Kallocain",
					Author:    "Karin Boye",
					Year:      uint16(1937),
					Available: false,
				},
				{
					ID:        4,
					Title:     "Kallocain",
					Author:    "Karin Boye",
					Year:      uint16(1937),
					Available: true,
				},
				{
					ID:        5,
					Title:     "Kallocain",
					Author:    "Karin Boye",
					Year:      uint16(1937),
					Available: false,
				},
				{
					ID:        6,
					Title:     "Kallocain",
					Author:    "Karin Boye",
					Year:      uint16(1937),
					Available: false,
				},
			},
		},
		{
			name:         "success, payload with available true",
			available:    true,
			availableNil: false,
			wantErr:      nil,
			wantBooks: []model.Book{
				{
					ID:        1,
					Title:     "1984",
					Author:    "George Orwell",
					Year:      uint16(1956),
					Available: true,
				},
				{
					ID:        4,
					Title:     "Kallocain",
					Author:    "Karin Boye",
					Year:      uint16(1937),
					Available: true,
				},
			},
		},
		{
			name:         "success, payload with available false",
			available:    false,
			availableNil: false,
			wantErr:      nil,
			wantBooks: []model.Book{
				{
					ID:        2,
					Title:     "Kallocain",
					Author:    "Karin Boye",
					Year:      uint16(1937),
					Available: false,
				},
				{
					ID:        3,
					Title:     "Kallocain",
					Author:    "Karin Boye",
					Year:      uint16(1937),
					Available: false,
				},
				{
					ID:        5,
					Title:     "Kallocain",
					Author:    "Karin Boye",
					Year:      uint16(1937),
					Available: false,
				},
				{
					ID:        6,
					Title:     "Kallocain",
					Author:    "Karin Boye",
					Year:      uint16(1937),
					Available: false,
				},
			},
		},
		{
			name:         "canceled",
			available:    false,
			availableNil: true,
			wantErr:      context.Canceled,
			wantCancel:   true,
			wantBooks:    []model.Book{},
		},
	}

	createTestRepo := func() *repository.BookRepository {
		return repository.CreateRepo(
			7,
			[]model.Book{
				{
					ID:        1,
					Title:     "1984",
					Author:    "George Orwell",
					Year:      uint16(1956),
					Available: true,
				},
				{
					ID:        2,
					Title:     "Kallocain",
					Author:    "Karin Boye",
					Year:      uint16(1937),
					Available: false,
				},
				{
					ID:        3,
					Title:     "Kallocain",
					Author:    "Karin Boye",
					Year:      uint16(1937),
					Available: false,
				},
				{
					ID:        4,
					Title:     "Kallocain",
					Author:    "Karin Boye",
					Year:      uint16(1937),
					Available: true,
				},
				{
					ID:        5,
					Title:     "Kallocain",
					Author:    "Karin Boye",
					Year:      uint16(1937),
					Available: false,
				},
				{
					ID:        6,
					Title:     "Kallocain",
					Author:    "Karin Boye",
					Year:      uint16(1937),
					Available: false,
				},
			},
		)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := createTestRepo()

			if repo == nil {
				t.Fatal("failed to create repository")
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if tt.wantCancel {
				cancel()
			}

			payload := repository.GetBooksPayload{}

			if !tt.availableNil {
				payload.Available = &tt.available
			}

			books, err := repo.GetBooks(ctx, payload)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected err %v but got %v", tt.wantErr, err)
			}

			if len(books) != len(tt.wantBooks) {
				t.Fatalf("expected len of books %d but got %d", len(tt.wantBooks), len(books))
			}

			for i := range tt.wantBooks {
				if books[i] != tt.wantBooks[i] {
					t.Fatalf("expected books %v \n but got %v", tt.wantBooks, books)
				}
			}

			for i := range tt.wantBooks {
				tt.wantBooks[i].ID++
			}

			if !tt.wantCancel {
				books, err = repo.GetBooks(ctx, payload)

				if err != nil {
					t.Fatalf("failed to get books")
				}

				for i := range tt.wantBooks {
					if tt.wantBooks[i] == books[i] {
						t.Fatal("repo chnaged expectedly")
					}
				}
			}
		})
	}
}
