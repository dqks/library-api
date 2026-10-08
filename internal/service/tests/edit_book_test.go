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
		name           string
		title          string
		author         string
		year           uint16
		available      bool
		cancel         bool
		id             int
		wantBook       model.Book
		wantErr        error
		repositoryBook model.Book
		repositoryErr  error
		titleNil       bool
		authorNil      bool
		yearNil        bool
		availableNil   bool
	}{
		{
			id:        1,
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
			repositoryBook: model.Book{
				ID:        1,
				Title:     "New Title",
				Author:    "New Author",
				Year:      uint16(2000),
				Available: true,
			},
			repositoryErr: nil,
		},
		{
			id:        1,
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
			id:        1,
			name:      "err, empty string author",
			title:     "New Title",
			author:    "",
			year:      uint16(0),
			available: true,
			cancel:    false,
			wantBook:  model.Book{},
			wantErr:   apperrors.ErrInvalidValues,
		},
		{
			id:        1,
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
			id:             1,
			name:           "err, canceled by service",
			title:          "New Title",
			author:         "New Author",
			year:           uint16(time.Now().Year()),
			available:      true,
			cancel:         true,
			wantBook:       model.Book{},
			wantErr:        context.Canceled,
			repositoryBook: model.Book{},
			repositoryErr:  context.Canceled,
		},
		{
			name:         "err, all fields nil",
			titleNil:     true,
			authorNil:    true,
			yearNil:      true,
			availableNil: true,
			wantBook:     model.Book{},
			wantErr:      apperrors.ErrRequiredFields,
		},
		{
			name:      "err, nil title",
			title:     "",
			author:    "New Author",
			year:      uint16(2000),
			available: true,
			cancel:    false,
			wantBook: model.Book{
				Title:     "New Title",
				Author:    "New Author",
				Year:      uint16(0),
				Available: true,
			},
			repositoryBook: model.Book{
				Title:     "New Title",
				Author:    "New Author",
				Year:      uint16(0),
				Available: true,
			},
			titleNil: true,
		},
		{
			name:      "err, nil author",
			title:     "New Title",
			author:    "",
			year:      uint16(2000),
			available: true,
			cancel:    false,
			wantBook: model.Book{
				Title:     "New Title",
				Author:    "New Author",
				Year:      uint16(0),
				Available: true,
			},
			repositoryBook: model.Book{
				Title:     "New Title",
				Author:    "New Author",
				Year:      uint16(0),
				Available: true,
			},
			authorNil: true,
		},
		{
			name:      "err, nil year",
			title:     "New Title",
			author:    "New Author",
			year:      uint16(0),
			available: true,
			cancel:    false,
			wantBook: model.Book{
				Title:     "New Title",
				Author:    "New Author",
				Year:      uint16(0),
				Available: true,
			},
			repositoryBook: model.Book{
				Title:     "New Title",
				Author:    "New Author",
				Year:      uint16(0),
				Available: true,
			},
			yearNil: true,
		},
		{
			name:      "err, nil available",
			title:     "New Title",
			author:    "New Author",
			year:      uint16(0),
			available: true,
			cancel:    false,
			wantBook: model.Book{
				Title:     "New Title",
				Author:    "New Author",
				Year:      uint16(0),
				Available: true,
			},
			repositoryBook: model.Book{
				Title:     "New Title",
				Author:    "New Author",
				Year:      uint16(0),
				Available: true,
			},
			availableNil: true,
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
				book: tt.repositoryBook,
				err:  tt.repositoryErr,
			})

			input := service.EditBookInput{}

			if !tt.authorNil {
				input.Author = &tt.author
			}

			if !tt.titleNil {
				input.Title = &tt.title
			}

			if !tt.yearNil {
				input.Year = &tt.year
			}

			if !tt.availableNil {
				input.Available = &tt.available
			}

			book, err := s.EditBookByID(ctx, tt.id, input)

			if book != tt.wantBook {
				t.Fatalf("expected book %v but got %v", tt.wantBook, book)
			}

			if err != tt.wantErr {
				t.Fatalf("expected err %v but got %v", tt.wantErr, err)
			}
		})
	}
}
