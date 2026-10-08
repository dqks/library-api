package tests

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/service"
	"testing"
	"time"
)

func TestCreateBook(t *testing.T) {
	tests := []struct {
		name           string
		title          string
		titleNil       bool
		author         string
		authorNil      bool
		year           uint16
		yearNil        bool
		available      bool
		availableNil   bool
		cancel         bool
		wantBook       model.Book
		wantErr        error
		repositoryBook model.Book
		repositoryErr  error
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
			repositoryBook: model.Book{
				ID:        1,
				Title:     "New Title",
				Author:    "New Author",
				Year:      uint16(2000),
				Available: true,
			},
			wantErr: nil,
		},
		{
			name:      "err, nil title",
			title:     "",
			author:    "New Author",
			year:      uint16(2000),
			available: true,
			cancel:    false,
			wantBook:  model.Book{},
			wantErr:   apperrors.ErrRequiredFields,
			titleNil:  true,
		},
		{
			name:      "err, nil author",
			title:     "New Title",
			author:    "",
			year:      uint16(2000),
			available: true,
			cancel:    false,
			wantBook:  model.Book{},
			wantErr:   apperrors.ErrRequiredFields,
			authorNil: true,
		},
		{
			name:      "err, nil year",
			title:     "New Title",
			author:    "",
			year:      uint16(0),
			available: true,
			cancel:    false,
			wantBook:  model.Book{},
			wantErr:   apperrors.ErrRequiredFields,
			yearNil:   true,
		},
		{
			name:         "err, nil available",
			title:        "New Title",
			author:       "",
			year:         uint16(0),
			available:    true,
			cancel:       false,
			wantBook:     model.Book{},
			wantErr:      apperrors.ErrRequiredFields,
			availableNil: true,
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
			title:     "New Title",
			author:    "",
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
			name:          "err, repository context canceled",
			title:         "New Title",
			author:        "New Author",
			year:          uint16(2000),
			available:     true,
			wantBook:      model.Book{},
			wantErr:       context.Canceled,
			repositoryErr: context.Canceled,
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

			input := service.CreateBookInput{}

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

			book, err := s.CreateBook(ctx, input)

			if err != tt.wantErr {
				t.Fatalf("expected err %v but got %v", tt.wantErr, err)
			}

			if book != tt.wantBook {
				t.Fatalf("expected book %v but got %v", tt.wantBook, book)
			}
		})
	}
}
