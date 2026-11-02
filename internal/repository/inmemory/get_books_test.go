package inmemory_test

import (
	"context"
	"library-api/internal/model"
	"library-api/internal/repository"
	"library-api/internal/repository/inmemory"
	"library-api/internal/testutil"
	"testing"
)

func TestGetBooks(t *testing.T) {
	tests := []struct {
		cancel       bool
		name         string
		available    bool
		availableNil bool
		wantErr      error
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
			cancel:       true,
			wantBooks:    []model.Book{},
		},
	}

	testBooks := []model.Book{
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
	}
	nextID := 7

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := inmemory.Create(
				nextID,
				testBooks,
			)

			if repo == nil {
				t.Fatal("failed to create repository")
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if tt.cancel {
				cancel()
			}

			payload := repository.GetBooksPayload{}

			if !tt.availableNil {
				payload.Available = &tt.available
			}

			books, err := repo.GetBooks(ctx, payload)

			fields := testutil.GetBooksFields{
				Err:       err,
				Books:     books,
				Cancel:    tt.cancel,
				Payload:   payload,
				Ctx:       ctx,
				Repo:      repo,
				WantErr:   tt.wantErr,
				WantBooks: tt.wantBooks,
			}

			testutil.CheckGetBooks(t, &fields)
		})
	}
}
