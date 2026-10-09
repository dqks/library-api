package tests

import (
	"context"
	"library-api/internal/apperrors"
	"library-api/internal/model"
	"library-api/internal/repository"
	"library-api/internal/service"
	"reflect"
	"testing"
	"time"
)

type editBookPayload struct {
	title     string
	author    string
	year      uint16
	available bool
}

func TestEditBookByID(t *testing.T) {
	tests := []struct {
		name            string
		title           string
		author          string
		year            uint16
		available       bool
		cancel          bool
		id              int
		wantBook        model.Book
		wantErr         error
		repositoryBook  model.Book
		repositoryErr   error
		titleNil        bool
		authorNil       bool
		yearNil         bool
		availableNil    bool
		receivedPayload editBookPayload
		getInRepo       bool
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
			receivedPayload: editBookPayload{
				title:     "New Title",
				author:    "New Author",
				year:      uint16(2000),
				available: true,
			},
			getInRepo: true,
		},
		{
			id:             1,
			name:           "err, not found",
			title:          "New Title",
			author:         "New Author",
			year:           uint16(2000),
			available:      true,
			cancel:         false,
			wantBook:       model.Book{},
			wantErr:        apperrors.ErrNotFound,
			repositoryBook: model.Book{},
			repositoryErr:  apperrors.ErrNotFound,
			receivedPayload: editBookPayload{
				title:     "New Title",
				author:    "New Author",
				year:      uint16(2000),
				available: true,
			},
			getInRepo: true,
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
			wantBook:       model.Book{},
			wantErr:        context.Canceled,
			cancel:         true,
			repositoryBook: model.Book{},
		},
		{
			id:             1,
			name:           "err, canceled by repository",
			title:          "New Title",
			author:         "New Author",
			year:           uint16(time.Now().Year()),
			available:      true,
			wantBook:       model.Book{},
			wantErr:        context.Canceled,
			repositoryBook: model.Book{},
			repositoryErr:  context.Canceled,
			receivedPayload: editBookPayload{
				title:     "New Title",
				author:    "New Author",
				year:      uint16(time.Now().Year()),
				available: true,
			},
			getInRepo: true,
		},
		{
			id:           1,
			name:         "err, all fields nil",
			titleNil:     true,
			authorNil:    true,
			yearNil:      true,
			availableNil: true,
			wantBook:     model.Book{},
			wantErr:      apperrors.ErrRequiredFields,
		},
		{
			id:        1,
			name:      "nil title",
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
			receivedPayload: editBookPayload{
				title:     "New Title",
				author:    "New Author",
				year:      uint16(2000),
				available: true,
			},
		},
		{
			id:        1,
			name:      "nil author",
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
			receivedPayload: editBookPayload{
				title:     "New Title",
				author:    "New Author",
				year:      uint16(2000),
				available: true,
			},
		},
		{
			id:        1,
			name:      "nil year",
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
			receivedPayload: editBookPayload{
				title:     "New Title",
				author:    "New Author",
				year:      uint16(2000),
				available: true,
			},
		},
		{
			id:        1,
			name:      "nil available",
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
			receivedPayload: editBookPayload{
				title:     "New Title",
				author:    "New Author",
				year:      uint16(0),
				available: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if tt.cancel {
				cancel()
			}

			repo := &MockRepo{
				book: tt.repositoryBook,
				err:  tt.repositoryErr,
			}

			s := service.Create(repo)

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

			if tt.getInRepo {
				if repo.gotInside {
					if tt.id != repo.id {
						t.Fatalf("expected to call repository with id %d but got %d", tt.id, repo.id)
					}

					expectedPayload := repository.EditBookPayload{}

					if !tt.authorNil {
						expectedPayload.Author = &tt.receivedPayload.author
					}

					if !tt.titleNil {
						expectedPayload.Title = &tt.receivedPayload.title
					}

					if !tt.yearNil {
						expectedPayload.Year = &tt.receivedPayload.year
					}

					if !tt.availableNil {
						expectedPayload.Available = &tt.receivedPayload.available
					}

					if !reflect.DeepEqual(repo.editBookPayload, expectedPayload) {
						t.Fatalf(
							"expected repository payload %#v\nbut got %#v",
							expectedPayload,
							repo.editBookPayload,
						)
					}
				} else {
					t.Fatalf("expected to get in repository but didn't")
				}
			}

			if book != tt.wantBook {
				t.Fatalf("expected book %v but got %v", tt.wantBook, book)
			}

			if err != tt.wantErr {
				t.Fatalf("expected err %v but got %v", tt.wantErr, err)
			}
		})
	}
}
