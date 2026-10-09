package tests

import (
	"context"
	"library-api/internal/model"
	"library-api/internal/repository"
	"library-api/internal/service"
	"reflect"
	"testing"
)

type getBooksPayload struct {
	available bool
}

func TestGetBooks(t *testing.T) {
	tests := []struct {
		name            string
		available       bool
		availableNil    bool
		cancel          bool
		wantBooks       []model.Book
		wantErr         error
		repositoryErr   error
		repositoryBooks []model.Book
		receivedPayload getBooksPayload
		getInRepo       bool
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
			wantErr:       nil,
			repositoryErr: nil,
			repositoryBooks: []model.Book{
				{
					ID:        1,
					Title:     "Mock Title",
					Author:    "Mock Author",
					Year:      2000,
					Available: false,
				},
			},
			receivedPayload: getBooksPayload{
				available: false,
			},
			getInRepo: true,
		},
		{
			name:         "success, available true",
			available:    true,
			availableNil: false,
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
			wantErr:       nil,
			repositoryErr: nil,
			repositoryBooks: []model.Book{
				{
					ID:        1,
					Title:     "Mock Title",
					Author:    "Mock Author",
					Year:      2000,
					Available: false,
				},
			},
			receivedPayload: getBooksPayload{
				available: true,
			},
			getInRepo: true,
		},
		{
			name:      "canceled by service",
			available: false,
			wantBooks: []model.Book{},
			wantErr:   context.Canceled,
			cancel:    true,
		},
		{
			name:            "canceled by repository",
			available:       false,
			wantBooks:       []model.Book{},
			wantErr:         context.Canceled,
			repositoryBooks: nil,
			repositoryErr:   context.Canceled,
			receivedPayload: getBooksPayload{
				available: false,
			},
			getInRepo: true,
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
				books: tt.repositoryBooks,
				err:   tt.repositoryErr,
			}

			s := service.Create(repo)

			params := service.GetBooksQueryParams{}

			if !tt.availableNil {
				params.Available = &tt.available
			}

			books, err := s.GetBooks(ctx, params)

			if tt.getInRepo {
				if repo.gotInside {
					payload := repository.GetBooksPayload{}

					if !tt.availableNil {
						payload.Available = &tt.receivedPayload.available
					}

					if !reflect.DeepEqual(payload, repo.getBooksPayload) {
						t.Fatalf("expected to call repository with payload %v\nbut got %v", payload, repo.getBooksPayload)
					}
				} else {
					t.Fatalf("expected to get in repository but didn't")
				}
			}

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
